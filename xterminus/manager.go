
package xterminus

import (
    "sync"
    "time"
    "log"
    "fmt"

    "chat.com/xcache"
    "chat.com/xutil"
    "chat.com/xmodel"
)

type Manager struct {
    Clients			map[*Client]bool			// all connection with peer websocket client
    Clocker		    sync.RWMutex				// read-write lock for Clients
    Users			map[string]*Client			// client as logined. key as appId_uuid
    Ulocker	    	sync.RWMutex				// read-write lock for Users
    Register		chan *Client				// while Client connect
    Login			chan *login					// while user logined
    Unregister		chan *Client				// while Client disconnect
    Broadcast 		chan []byte                 // broadcast the message to all Users
}

func NewManager() (*Manager) {
    return &Manager{
        Clients:    make(map[*Client]bool),
        Users:      make(map[string]*Client),
        Register:   make(chan *Client, 1024),
        Login:      make(chan *login, 1024),
        Unregister: make(chan *Client, 1024),
        Broadcast:  make(chan []byte, 1024),
    }
}

func getUserKey(appId uint32, userId string) (string) {
    return fmt.Sprintf("%d_%s", appId, userId)
}

func (this *Manager) clientsRange(f func(*Client, bool) (bool)) {
    this.Clocker.RLock()
    defer this.Clocker.RUnlock()

    for c, b := range this.Clients {
        res := f(c, b)
        if res == false {
            return
        }
    }
    return
}

// 连接是否存在
func (this *Manager) existClient(c *Client) (bool) {
    this.Clocker.RLock()
    defer this.Clocker.RUnlock()
    _, exists := this.Clients[c]

    return exists
}

// 获取所有连接
func (this *Manager) getClients() (map[*Client]bool) {
    this.Clocker.RLock()
    defer this.Clocker.RUnlock()

    return this.Clients
}

func (this *Manager) getClientsLen() (int) {
    this.Clocker.RLock()
    defer this.Clocker.RUnlock()

    return len(this.Clients)
}

func (this *Manager) addClient(c *Client) {
    this.Clocker.Lock()
    defer this.Clocker.Unlock()

    this.Clients[c] = true
}

func (this *Manager) delClient(c *Client) {
    this.Clocker.Lock()
    defer this.Clocker.Unlock()

    if _, exists := this.Clients[c]; exists {
        delete(this.Clients, c)
    }
}

// acquire Client connection specified by appId and userId
func (this *Manager) GetUserClient(appId uint32, userId string) (*Client) {
    this.Ulocker.RLock()
    defer this.Ulocker.RUnlock()

    uk := getUserKey(appId, userId)
    if c, exists := this.Users[uk]; exists {
        return c
    }

    return nil
}

func (this *Manager) getUsersLen() (int) {
    this.Ulocker.RLock()
    defer this.Ulocker.RUnlock()

    return len(this.Users)
}

func (this *Manager) addUser(uk string, c *Client) {
    this.Ulocker.Lock()
    defer this.Ulocker.Unlock()

    this.Users[uk] = c
}

func (this *Manager) delUser(c *Client) (bool) {
    this.Ulocker.Lock()
    defer this.Ulocker.Unlock()

    uk := getUserKey(c.AppId, c.UserId)
    bl := false
    if v, exists := this.Users[uk]; exists {
        // is the same machine?
        if v.Addr != c.Addr {
            return bl
        }
        delete(this.Users, uk)
        bl = true
    }

    return bl
}

func (this *Manager) getUserKeys() ([]string) {
    this.Ulocker.RLock()
    defer this.Ulocker.RUnlock()

    uks := make([]string, 0)
    for k := range this.Users {
        uks = append(uks, k)
    }

    return uks
}

// get all user's id scoped by appId
func (this *Manager) GetUserList(appId uint32) ([]string) {
    this.Ulocker.RLock()
    defer this.Ulocker.RUnlock()

    uls := make([]string, 0)
    for _, c := range this.Users {
        if c.AppId == appId {
            uls = append(uls, c.UserId)
        }
    }

    log.Printf("-DBUG- users count: %d", len(this.Users))

    return uls
}

// get all client connections
func (this *Manager) GetUserClients() ([]*Client) {
    this.Ulocker.RLock()
    defer this.Ulocker.RUnlock()

    cs := make([]*Client, 0)
    for _, c := range this.Users {
        cs = append(cs, c)
    }

    return cs
}

func (this *Manager) sendMessageToOther(message []byte, appId uint32, ignoreClient *Client) {
    cs := this.GetUserClients()
    for _, c := range cs {
        if c != ignoreClient && c.AppId == appId {
            c.sendTo(message)
        }
    }
}

// event: connect
func (this *Manager) EventRegister(c *Client) {
    this.addClient(c)

    log.Printf("-INFO- client(%s) connected",  c.Addr)
}

// event: login
func (this *Manager) EventLogin(l *login) {
    c := l.Client

    if this.existClient(c) {
        uk := l.GetKey()
        this.addUser(uk, c)
    }

    log.Printf("-INFO- client(%d_%s) been managed", l.AppId, l.UserId)

    orderId := xutil.GetOrderIdTime()
    SendMessageToAll(l.AppId, l.UserId, orderId, xmodel.MessageCmdEnter, "哈喽~")
    
}

func (this *Manager) EventUnregister(c *Client) {
    this.delClient(c)

    if ! this.delUser(c) {
        return
    }

    userOnline, err := xcache.GetUserOnlineInfo(c.GetKey())
    if err == nil {
        userOnline.Logout()
        xcache.SetUserOnlineInfo(c.GetKey(), userOnline)
    }

    log.Printf("-INFO- client(%d_%s) disconnected", c.AppId, c.UserId)

    if c.UserId != "" {
        orderId := xutil.GetOrderIdTime()
        SendMessageToAll(c.AppId, c.UserId, orderId, xmodel.MessageCmdExit, "用户已经离开~")
    }
}

func (this *Manager) Start() {
    for {
        select {
        case c := <- this.Register:
            this.EventRegister(c)
        case l := <- this.Login:
            this.EventLogin(l)
        case c := <- this.Unregister:
            this.EventUnregister(c)
        case msg := <- this.Broadcast:
            cs := this.getClients()
            for c := range cs {
                select {
                case c.Send <- msg:
                default:
                    close(c.Send)
                }
            }
        }
    }
}

func GetManagerInfo(isDebug string) (map[string]interface{}) {
    mi := make(map[string]interface{})

    mi["clientsLen"] = _manager.getClientsLen()
    mi["usersLen"] = _manager.getUsersLen()
    mi["chanRegisterLen"] = len(_manager.Register)
    mi["chanLoginLen"] = len(_manager.Login)
    mi["chanUnregisterLen"] = len(_manager.Unregister)
    mi["chanBroadcastLen"] = len(_manager.Broadcast)

    if isDebug == "true" {
        addrList := make([]string, 0)
        _manager.clientsRange(func(c *Client, v bool)(bool) {
            addrList = append(addrList, c.Addr)
            return true
        })

        users := _manager.getUserKeys()
        mi["clients"] = addrList        // Client List
        mi["users"] = users             // User List
    }

    return mi
}

func GetUserClient(appId uint32, userId string) (*Client) {
    return _manager.GetUserClient(appId, userId)
}

func CleanTimeoutConnections() {
    currentTime := uint64(time.Now().Unix())

    cs := _manager.getClients()
    for c := range cs {
        if c.IsHeartbeatTimeout(currentTime) {
            c.Socket.Close()
            log.Printf("-ERRO- client(%d_%s loginTime=%d, heartbeatTime=%d) connected timeout", c.AppId, c.UserId, c.LoginTime, c.HeartbeatTime)
        }
    }
}

func getUserList(appId uint32) ([]string) {
    return _manager.GetUserList(appId)
}

func sendMessageToAppId(appId uint32, userId string, data string) {
    ignoreClient := _manager.GetUserClient(appId, userId)
    _manager.sendMessageToOther([]byte(data), appId, ignoreClient)
}
