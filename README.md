# Go's Decoder Ring: Cloud based Secrets Management Made Easy

```go
	var s Secrets
	l := NewLoader(&s)
	l.With(GoogleCloudLoader)
	err := l.Load()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
```