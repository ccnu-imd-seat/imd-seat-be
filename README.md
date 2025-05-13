# imd-seat-be
信息管理学院选座系统后端

## 目录结构说明
```
.
├── etc/            # 配置文件
│   └── config.yaml # 应用配置
├── internal/
│   ├── config/     # 配置结构体定义
│   ├── handler/    # 路由处理层
│   ├── logic/      # 业务逻辑层  
│   ├── model/      # 数据模型层
│   ├── svc/        # 服务上下文
│   ├── types/      # 请求/响应DTO
|   └── sql/        # 数据表设计
├── seat.api        # API定义文件
└── seat.go         # 主入口
```

## API开发流程
1. 在`seat.api`中定义接口
```go
type (
    ReserveRequest {
        SeatID int64 `json:"seatId"`
        Date   string `json:"date"`
    }
    
    ReserveResponse {
        ReservationID int64 `json:"reservationId"`
    }
)

service imd-seat-be {
    @handler ReserveSeatHandler
    post /api/reserve (ReserveRequest) returns (ReserveResponse)
}
```

2. 生成代码
```bash
goctl api go -api seat.api -dir .
```

3. 实现业务逻辑
```go
// internal/logic/reserveseatlogic.go
func (l *ReserveSeatLogic) ReserveSeat(req *types.ReserveRequest) (*types.ReserveResponse, error) {
    // 业务逻辑实现
    return &types.ReserveResponse{ReservationID: 123}, nil
}
```

## 数据库设计
主要表结构见`internal/sql/data.sql`，包含：
- 座位表(seats)  
- 预约表(reservations)

## 接口设计
[接口设计](https://apifox.com/apidoc/shared/2723d5c5-2c20-467d-a3d6-5138f50b0e4b)