// Code generated by "genprotocol -ver=de39b3963afaad4ed2557809c9208a5676bf9321c0823e0dac15bc4db6d51552 -basedir=. -prefix=w3d -statstype=int"

package w3d_handlersp

/* bytes base demux fn map template

var DemuxRsp2BytesFnMap = [...]func(me interface{}, hd w3d_packet.Header, rbody []byte) error {
w3d_idcmd.Invalid : bytesRecvRspFn_Invalid,
w3d_idcmd.Login : bytesRecvRspFn_Login,
w3d_idcmd.Heartbeat : bytesRecvRspFn_Heartbeat,
w3d_idcmd.Chat : bytesRecvRspFn_Chat,

}

	func bytesRecvRspFn_Invalid(me interface{}, hd w3d_packet.Header, rbody []byte) error {
		robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*w3d_obj.RspInvalid_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

	func bytesRecvRspFn_Login(me interface{}, hd w3d_packet.Header, rbody []byte) error {
		robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*w3d_obj.RspLogin_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

	func bytesRecvRspFn_Heartbeat(me interface{}, hd w3d_packet.Header, rbody []byte) error {
		robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*w3d_obj.RspHeartbeat_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

	func bytesRecvRspFn_Chat(me interface{}, hd w3d_packet.Header, rbody []byte) error {
		robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*w3d_obj.RspChat_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

*/
