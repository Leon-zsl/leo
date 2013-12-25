/* this is rpc client
 */

package base

import (
	"net/rpc"
	"strconv"
	"strings"
)

//safe for goroutine
type RpcCallback interface {
	HandlerReplay(reply interface{}, err error)
}

type RpcClient struct {
	running bool
	ip      string
	port    int
	cl      *rpc.Client
}

func NewRpcClient(ip string, port int) (client *RpcClient, err error) {
	client = new(RpcClient)
	err = client.init(ip, port)
	return
}

func (client *RpcClient) init(ip string, port int) error {
	client.ip = ip
	client.port = port
	return nil
}

func (client *RpcClient) Start() error {
	arr := []string{client.ip, strconv.Itoa(client.port)}
	addr := strings.Join(arr, ":")

	cl, err := rpc.Dial("tcp", addr)
	if err != nil {
		return err
	}

	client.cl = cl
	client.running = true
	return nil
}

func (client *RpcClient) Close() {
	client.running = false
	if client.cl != nil {
		client.cl.Close()
		client.cl = nil
	}
}

func (client *RpcClient) IP() string {
	return client.ip
}

func (client *RpcClient) Port() int {
	return client.port
}

func (client *RpcClient) Call(method string, args interface{}, reply interface{}) error {
	return client.cl.Call(method, args, reply)
}

func (client *RpcClient) Go(method string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call {
	return client.cl.Go(method, args, reply, done)
}

func (client *RpcClient) CallAsync(method string, args interface{}, reply interface{}, cb RpcCallback) {
	go client.rpc_call(method, args, reply, cb)
}

func (client *RpcClient) rpc_call(method string, args interface{}, reply interface{}, cb RpcCallback) {
	call := client.cl.Go(method, args, reply, nil)
	rspcall := <-call.Done
	if cb != nil {
		cb.HandlerReplay(rspcall.Reply, rspcall.Error)
	} else if rspcall.Error != nil {
		LoggerIns.Error("rpc async error", rspcall.Error)
	}
}
