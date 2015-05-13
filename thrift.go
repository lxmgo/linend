package revel

import (
	"crypto/tls"
	"git.apache.org/thrift.git/lib/go/thrift"
	"strconv"
)

var processor *thrift.TMultiplexedProcessor

type Handler struct{
	
}

//Start thrift server while init
func init() {
	processor =thrift.NewTMultiplexedProcessor() 
}

//Register processer, called in app code
func RegisterProcessor(name string, p thrift.TProcessor){
	processor.RegisterProcessor(name,p)
	TRACE.Printf("Registered RegisterProcessor: %s", name)
}

func RunThrift(){
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTTransportFactory()
	
	addr := ThriftAddr + ":" + strconv.Itoa(ThriftPort)

	if err := runServer(processor,transportFactory, protocolFactory, addr, ThriftSecure); err != nil {
		ERROR.Fatal("Run thrift error ",err)
	}
}
func runServer(tMultiplexedProcessor *thrift.TMultiplexedProcessor,transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) error {
	var transport thrift.TServerTransport
	var err error
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(addr)
	}

	if err != nil {
		return err
	}

	server := thrift.NewTSimpleServer4(tMultiplexedProcessor, transport, transportFactory, protocolFactory)

	INFO.Println("Running thrift server on ", addr)
	return server.Serve()
}