# imd-seat-be
信息管理学院选座系统后端

## 目录结构说明
```
.
├── api/            # API接口定义文件
│   ├── index.api   # 代码生成入口API
│   └── *.api       # 各模块API定义
├── etc/            # 配置文件
│   └── config.yaml # 应用配置
├── internal/
│   ├── config/     # 配置结构体定义
│   ├── handler/    # 路由处理层
│   ├── logic/      # 业务逻辑层
│   ├── model/      # 数据模型层
│   ├── svc/        # 服务上下文
│   ├── types/      # 请求/响应DTO
│   └── sql/        # 数据表设计
├── go.mod          # Go模块定义
├── go.sum          # 依赖校验
├── main.go         # 主入口
└── Makefile        # 构建脚本
```

## API开发流程
1. 在对应的`.api`文件中定义接口(如`reservation.api`)
```go
@server(
    prefix: /api/reservation
    group: reservation
)
service imd-seat-be {
    @handler ReserveSeatHandler
    post /reserve (ReserveRequest) returns (ReserveResponse)
    
    @handler CancelReservationHandler  
    post /cancel (CancelRequest) returns (CancelResponse)
}

type (
    ReserveRequest {
        SeatID int64 `json:"seatId"`
        Date   string `json:"date"`
    }
    
    ReserveResponse {
        ReservationID int64 `json:"reservationId"`
    }
    
    CancelRequest {
        ReservationID int64 `json:"reservationId"`
    }
    
    CancelResponse {
        Success bool `json:"success"`
    }
)
```

2. 生成代码
```bash
make api
```

3. 实现业务逻辑
```go
// internal/logic/reserveseatlogic.go
func (l *ReserveSeatLogic) ReserveSeat(req *types.ReserveRequest) (*types.ReserveResponse, error) {
    // 1. 参数校验
    // 2. 检查座位可用性
    // 3. 创建预约记录
    return &types.ReserveResponse{ReservationID: newID}, nil
}
```

## 数据库设计
主要表结构见`internal/sql/data.sql`，包含：
- 座位表(seats)  
- 预约表(reservations)

## 接口设计
[接口设计](https://apifox.com/apidoc/shared/2723d5c5-2c20-467d-a3d6-5138f50b0e4b)
