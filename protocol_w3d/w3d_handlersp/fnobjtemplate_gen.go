// Code generated by "genprotocol -ver=f5b1d289172cf84ad5d01b91533408be6b17961cf28ddd6fe767224298a8aedd -basedir=. -prefix=w3d -statstype=int"

package w3d_handlersp

/* obj base demux fn map template

var DemuxRsp2ObjFnMap = [...]func(me interface{}, hd w3d_packet.Header, body interface{}) error {
w3d_idcmd.Invalid : objRecvRspFn_Invalid,
w3d_idcmd.MakeTeam : objRecvRspFn_MakeTeam,
w3d_idcmd.Act : objRecvRspFn_Act,
w3d_idcmd.Heartbeat : objRecvRspFn_Heartbeat,

}

	func objRecvRspFn_Invalid(me interface{}, hd w3d_packet.Header, body interface{}) error {
		robj , ok := body.(*w3d_obj.RspInvalid_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

	func objRecvRspFn_MakeTeam(me interface{}, hd w3d_packet.Header, body interface{}) error {
		robj , ok := body.(*w3d_obj.RspMakeTeam_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

	func objRecvRspFn_Act(me interface{}, hd w3d_packet.Header, body interface{}) error {
		robj , ok := body.(*w3d_obj.RspAct_data)
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

*/
