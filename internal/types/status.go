package types

//预约状态的状态码
const (
	BookedStatus     = "booked"      //已预约
	EffectiveStatus  = "effective"   //已生效未签到
	CheckedInStatus  = "checked in"  //已签到
	CanceldStatus    = "canceld"     //已取消
	CompletedStatus  = "completed"   //预约已完成
	OneDefaultStatus = "one default" //已未到一次
	ViolatedStatus   = "violated"    // 已违约
	AvaliableStatus  = "avaliable"   //可预约
)
