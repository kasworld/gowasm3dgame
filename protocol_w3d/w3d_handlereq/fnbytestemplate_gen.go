// Code generated by "genprotocol -ver=213afa194ef0e682076c6a0cbf801946c13d343cc54330be7c4557e46057a498 -basedir=. -prefix=w3d -statstype=int"

package w3d_handlereq

/* bytes base fn map api template , unmarshal in api
	var DemuxReq2BytesAPIFnMap = [...]func(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error){
	w3d_idcmd.Invalid: bytesAPIFn_ReqInvalid,
w3d_idcmd.EnterStage: bytesAPIFn_ReqEnterStage,
w3d_idcmd.ChatToStage: bytesAPIFn_ReqChatToStage,
w3d_idcmd.Heartbeat: bytesAPIFn_ReqHeartbeat,

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

	func bytesAPIFn_ReqChatToStage(
		me interface{}, hd w3d_packet.Header, rbody []byte) (
		w3d_packet.Header, interface{}, error) {
		// robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*w3d_obj.ReqChatToStage_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := w3d_packet.Header{
			ErrorCode : w3d_error.None,
		}
		sendBody := &w3d_obj.RspChatToStage_data{
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

*/
