// Code generated by "genprotocol -ver=de39b3963afaad4ed2557809c9208a5676bf9321c0823e0dac15bc4db6d51552 -basedir=. -prefix=w3d -statstype=int"

package w3d_handlersp

/* obj base demux fn map template

var DemuxRsp2ObjFnMap = [...]func(me interface{}, hd w3d_packet.Header, body interface{}) error {
w3d_idcmd.Invalid : objRecvRspFn_Invalid,
w3d_idcmd.Login : objRecvRspFn_Login,
w3d_idcmd.Heartbeat : objRecvRspFn_Heartbeat,
w3d_idcmd.Chat : objRecvRspFn_Chat,

}

	func objRecvRspFn_Invalid(me interface{}, hd w3d_packet.Header, body interface{}) error {
		robj , ok := body.(*w3d_obj.RspInvalid_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

	func objRecvRspFn_Login(me interface{}, hd w3d_packet.Header, body interface{}) error {
		robj , ok := body.(*w3d_obj.RspLogin_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

	func objRecvRspFn_Heartbeat(me interface{}, hd w3d_packet.Header, body interface{}) error {
		robj , ok := body.(*w3d_obj.RspHeartbeat_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

	func objRecvRspFn_Chat(me interface{}, hd w3d_packet.Header, body interface{}) error {
		robj , ok := body.(*w3d_obj.RspChat_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

*/
