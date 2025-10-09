package main
import(
	"net/http"
	"github.com/max-durnea/ByteBucket/internal/auth"
)

func (cfg *apiConfig) JwtMiddleware(next http.Handler) http.Handler{
	return http.HandleFunc(func (w http.ResponseWriter, r *http.Request){
		//extract token from the header
		token,err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w,404,"Access Denied")
		}
		user_id, err := auth.ValidateJWT(token,cfg.tokenSecret)
		if err != nil {
			respondWithError(w,404,"Access Denied")
		}
		//store the id in the context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}