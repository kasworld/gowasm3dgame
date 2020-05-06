// Code generated by "genprotocol -ver=fd815e2fbc449528b4fb5d55480af0a03b4bfaf074ff2c5570d2e5a3ce03896b -basedir=. -prefix=w3d -statstype=int"

package w3d_handlereq

/* bytes base fn map api template , unmarshal in api
	var DemuxReq2BytesAPIFnMap = [...]func(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error){
	w3d_idcmd.Invalid: bytesAPIFn_ReqInvalid,
w3d_idcmd.Login: bytesAPIFn_ReqLogin,
w3d_idcmd.Heartbeat: bytesAPIFn_ReqHeartbeat,
w3d_idcmd.Chat: bytesAPIFn_ReqChat,
w3d_idcmd.EnterStage: bytesAPIFn_ReqEnterStage,

}   // DemuxReq2BytesAPIFnMap

	func bytesAPIFn_ReqInvalid(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error) {
		// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*w3d_obj.ReqInvalid_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := w3d_packet.Header{
			ErrorCode : w3d_error.None,
		}
		sendBody := &w3d_obj.RspInvalid_data{
		}
		return sendHeader, sendBody, nil
	}

	func bytesAPIFn_ReqLogin(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error) {
		// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*w3d_obj.ReqLogin_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := w3d_packet.Header{
			ErrorCode : w3d_error.None,
		}
		sendBody := &w3d_obj.RspLogin_data{
		}
		return sendHeader, sendBody, nil
	}

	func bytesAPIFn_ReqHeartbeat(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error) {
		// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*w3d_obj.ReqHeartbeat_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := w3d_packet.Header{
			ErrorCode : w3d_error.None,
		}
		sendBody := &w3d_obj.RspHeartbeat_data{
		}
		return sendHeader, sendBody, nil
	}

	func bytesAPIFn_ReqChat(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error) {
		// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*w3d_obj.ReqChat_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := w3d_packet.Header{
			ErrorCode : w3d_error.None,
		}
		sendBody := &w3d_obj.RspChat_data{
		}
		return sendHeader, sendBody, nil
	}

	func bytesAPIFn_ReqEnterStage(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error) {
		// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*w3d_obj.ReqEnterStage_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := w3d_packet.Header{
			ErrorCode : w3d_error.None,
		}
		sendBody := &w3d_obj.RspEnterStage_data{
		}
		return sendHeader, sendBody, nil
	}

*/
