package userman;
// pxnMetrics Broker - rpc user manager

import(
	Context  "context"
	GRPC     "google.golang.org/grpc"
	GStatus  "google.golang.org/grpc/status"
	GCodes   "google.golang.org/grpc/codes"
	UtilsRPC "github.com/PoiXson/pxnGoCommon/rpc"
	Configs  "github.com/PoiXson/pxnMetrics/broker/configs"
);



const KeyUserRPC = "user-rpc";



func NewUserManagerInterceptor(config *Configs.CfgBroker) GRPC.UnaryServerInterceptor {
	return func(ctx Context.Context, req any, info *GRPC.UnaryServerInfo,
			handler GRPC.UnaryHandler) (any, error) {
		username, ok := ctx.Value(UtilsRPC.KeyUsername).(string);
		if !ok || username == "" {
			return nil, GStatus.Errorf(
				GCodes.PermissionDenied,
				"Unable to find username",
			);
		}
		user, ok := config.Users[username];
		if !ok {
			return nil, GStatus.Errorf(
				GCodes.PermissionDenied,
				"Unable to find RPC user info",
			);
		}
		ctx = Context.WithValue(ctx, KeyUserRPC, &user);
		return handler(ctx, req);
	}
}
