package net

import (
	"fmt"
	getty "github.com/apache/dubbo-getty"
	gxnet "github.com/dubbogo/gost/net"
	gxsync "github.com/dubbogo/gost/sync"
	"net"
	"time"
	"xdago/config"
	"xdago/log"
)

const CronPeriod = 20e9

type XdagServer struct {
	server   getty.Server
	taskPool gxsync.GenericTaskPool
	config   *config.Config
}

func (s *XdagServer) Start(config *config.Config) {
	s.config = config
	addr := gxnet.HostAddress(config.NodeIp(), config.NodePort())
	serverOpts := []getty.ServerOption{getty.WithLocalAddress(addr)}
	s.taskPool = gxsync.NewTaskPoolSimple(128) // TODO: put task pool size into config file
	serverOpts = append(serverOpts, getty.WithServerTaskPool(s.taskPool))
	s.server = getty.NewTCPServer(serverOpts...)
	log.Debug("Listening for incoming connections", log.Ctx{"ip": config.NodeIp(), "port": config.NodePort()})
	s.server.RunEventLoop(s.newSession)
}

func (s *XdagServer) Close() {
	log.Debug("Closing XdagServer...")
	s.server.Close()
	s.taskPool.Close()
	log.Debug("XdagServer closed.")

}

func (s *XdagServer) newSession(session getty.Session) (err error) {
	tcpConn, ok := session.Conn().(*net.TCPConn)
	if !ok {
		panic(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection\n", session.Stat(), session.Conn()))
	}

	if err = tcpConn.SetNoDelay(true); err != nil {
		return err
	}
	if err = tcpConn.SetKeepAlive(true); err != nil {
		return err
	}
	if err = tcpConn.SetKeepAlivePeriod(10 * time.Second); err != nil {
		return err
	}
	if err = tcpConn.SetReadBuffer(262144); err != nil {
		return err
	}
	if err = tcpConn.SetWriteBuffer(524288); err != nil {
		return err
	}

	session.SetName("XdagServer")
	session.SetMaxMsgLen(128 * 1024) // max message package length is 128k
	session.SetReadTimeout(time.Duration(s.config.ConnectionReadTimeout()) * time.Millisecond)
	session.SetWriteTimeout(5 * time.Second)
	session.SetCronPeriod(int(CronPeriod / 1e6))
	session.SetWaitTime(time.Second)

	//session.SetPkgHandler(pkgHandler)
	//session.SetEventListener(EventListener)
	return nil
}
