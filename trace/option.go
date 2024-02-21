package trace

const TraceName = "lori"

type Options struct {
	Name     string  `json:"name"`
	Endpoint string  `json:"endpoint"`
	Sampler  float64 `json:"sampler"`
	Batcher  string  `json:"batcher"`
}

/*
	Name:     "服务名",
	Endpoint: "http://127.0.0.1:14268/api/traces",
	Sampler:  1.0,
	Batcher:  "jaeger",
*/
