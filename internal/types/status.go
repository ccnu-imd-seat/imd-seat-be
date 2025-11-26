package types

//预约状态的状态码
const (
	BookedStatus    = "booked"    //已预约
	EffectiveStatus = "effective" //已生效
	CancelledStatus = "cancelled" //已取消
	CompletedStatus = "completed" //已完成
	ViolatedStatus  = "violated"  //已违约
	AvaliableStatus = "available" //可预约
)
