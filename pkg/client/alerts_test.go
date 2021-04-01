package client_test
import (
  "grafana/pkg/client"
  "testing"
)
func TestAdd(t *testing.T){
  ret = mymath.Add(2,3)
  if ret != 5{
   t.Error("Expected 5, got",ret)
  }
}
func BenchmarkAdd(b *testing.B){
  for i := 0; i< b.N; i++{
    fmt.Sprintf("hello")
  }
}
