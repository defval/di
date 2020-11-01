package di

//
//import (
//	"fmt"
//	"net"
//	"net/http"
//	"reflect"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//)
//
//var (
//	s                = newDefaultSchema()
//	httpServerType   = reflect.TypeOf(&http.Server{})
//	httpServeMuxType = reflect.TypeOf(&http.ServeMux{})
//	netAddrType      = reflect.TypeOf(new(net.Addr)).Elem()
//)
//
//func provideHTTPServer(s schema) (_ reflect.Value, cleanup func(), err error) {
//	return reflect.ValueOf(&http.Server{}), nil, nil
//}
//
//func provideHTTPServeMux(s schema) (_ reflect.Value, cleanup func(), err error) {
//	return reflect.ValueOf(&http.ServeMux{}), nil, nil
//}
//
//func provideTCPAddr(s schema) (_ reflect.Value, cleanup func(), err error) {
//	return reflect.ValueOf(&net.TCPAddr{}), nil, nil
//}
//
//func provideUDPAddr(s schema) (_ reflect.Value, cleanup func(), err error) {
//	return reflect.ValueOf(&net.UDPAddr{}), nil, nil
//}
//
//func init() {
//	s.register(&node{rt: httpServerType, provide: provideHTTPServer})
//	s.register(&node{rt: httpServeMuxType, provide: provideHTTPServeMux})
//	s.register(&node{rt: netAddrType, provide: provideTCPAddr, tags: Tags{"type": "tcp", "md": "true"}})
//	s.register(&node{rt: netAddrType, provide: provideUDPAddr, tags: Tags{"type": "udp", "md": "true"}})
//}
//
//func Test_SchemaFind(t *testing.T) {
//	t.Run("find not existing type causes not found error", func(t *testing.T) {
//		fn, err := s.find(reflect.TypeOf(new(http.Request)), Tags{})
//		require.Nil(t, fn)
//		require.EqualError(t, err, "type *http.Request not exists in the container")
//	})
//
//	t.Run("find type without tags working correctly", func(t *testing.T) {
//		node, err := s.find(httpServerType, Tags{})
//		require.NoError(t, err)
//		require.Equal(t, fmt.Sprintf("%p", provideHTTPServer), fmt.Sprintf("%p", node.provide))
//		node, err = s.find(httpServeMuxType, Tags{})
//		require.NoError(t, err)
//		require.Equal(t, fmt.Sprintf("%p", provideHTTPServeMux), fmt.Sprintf("%p", node.provide))
//	})
//
//	t.Run("find type with tag working correctly", func(t *testing.T) {
//		node, err := s.find(netAddrType, Tags{"type": "tcp"})
//		require.NoError(t, err)
//		require.NotEqual(t, fmt.Sprintf("%p", provideUDPAddr), fmt.Sprintf("%p", node.provide))
//		require.Equal(t, fmt.Sprintf("%p", provideTCPAddr), fmt.Sprintf("%p", node.provide))
//
//		node, err = s.find(netAddrType, Tags{"type": "udp"})
//		require.NoError(t, err)
//		require.NotEqual(t, fmt.Sprintf("%p", provideTCPAddr), fmt.Sprintf("%p", node.provide))
//		require.Equal(t, fmt.Sprintf("%p", provideUDPAddr), fmt.Sprintf("%p", node.provide))
//	})
//
//	t.Run("tags don't fitting for a type causes a not exists error", func(t *testing.T) {
//		node, err := s.find(netAddrType, Tags{"type": "not_exists"})
//		require.Nil(t, node)
//		require.EqualError(t, err, "net.Addr[type:not_exists] not exists")
//	})
//
//	t.Run("trying to resolve multiple defined type without a group causes an error", func(t *testing.T) {
//		node, err := s.find(netAddrType, Tags{})
//		require.Nil(t, node)
//		require.EqualError(t, err, "multiple definitions of net.Addr, maybe you need to use group type: []net.Addr")
//	})
//
//	t.Run("trying to resolve multiple defined type with the same tag causes an error", func(t *testing.T) {
//		node, err := s.find(netAddrType, Tags{"md": "true"})
//		require.Nil(t, node)
//		require.EqualError(t, err, "multiple definitions of net.Addr[md:true], maybe you need to use group type: []net.Addr[md:true]")
//	})
//}
