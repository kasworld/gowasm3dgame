// Code generated by "genprotocol -ver=f5b1d289172cf84ad5d01b91533408be6b17961cf28ddd6fe767224298a8aedd -basedir=. -prefix=w3d -statstype=int"

package w3d_handlenoti

/* bytes base demux fn map template

var DemuxNoti2ByteFnMap = [...]func(me interface{}, hd w3d_packet.Header, rbody []byte) error {
w3d_idnoti.Invalid : bytesRecvNotiFn_Invalid,
w3d_idnoti.StageInfo : bytesRecvNotiFn_StageInfo,
w3d_idnoti.StatsInfo : bytesRecvNotiFn_StatsInfo,

}

	func bytesRecvNotiFn_Invalid(me interface{}, hd w3d_packet.Header, rbody []byte) error {
		robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*w3d_obj.NotiInvalid_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

	func bytesRecvNotiFn_StageInfo(me interface{}, hd w3d_packet.Header, rbody []byte) error {
		robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*w3d_obj.NotiStageInfo_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

	func bytesRecvNotiFn_StatsInfo(me interface{}, hd w3d_packet.Header, rbody []byte) error {
		robj, err := w3d_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*w3d_obj.NotiStatsInfo_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

*/