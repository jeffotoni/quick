package glog

import (
	"context"
	"fmt"
	"time"
)

// This function is named ExampleCreateCtx_basic()
// it with the Examples type.
func ExampleCreateCtx_basic() {
	ctx, cancel := CreateCtx().
		Set("TraceID", "abc-123").
		Build()
	defer cancel()

	traceID := GetCtx(ctx, "TraceID")
	fmt.Println("Trace ID:", traceID)

	// Output:
	// Trace ID: abc-123
}

// This function is named ExampleCreateCtx_multipleFields()
// it with the Examples type.
func ExampleCreateCtx_multipleFields() {
	ctx, cancel := CreateCtx().
		Set("TraceID", "xyz-789").
		Set("UserID", "42").
		Build()
	defer cancel()

	fields := GetCtxAll(ctx)
	fmt.Println(fields["TraceID"])
	fmt.Println(fields["UserID"])

	// Output:
	// xyz-789
	// 42
}

// This function is named ExampleCreateCtx_customTimeout()
// it with the Examples type.
func ExampleCreateCtx_customTimeout() {
	start := time.Now()

	ctx, cancel := CreateCtx().
		Set("Key", "val").
		Timeout(2 * time.Second).
		Build()
	defer cancel()

	deadline, _ := ctx.Deadline()
	fmt.Println("Timeout in ~2s:", deadline.Sub(start) < 3*time.Second)

	// Output:
	// Timeout in ~2s: true
}

// This function is named ExampleGetCtx_emptyOrNil()
// it with the Examples type.
func ExampleGetCtx_emptyOrNil() {
	fmt.Println(GetCtx(nil, "TraceID"))
	fmt.Println(GetCtx(context.Background()))
	fmt.Println(GetCtx(context.Background(), ""))

	// Output:
	//

}

// This function is named ExampleCtxBuilder_Set()
// it with the Examples type.
func ExampleCtxBuilder_Set() {
	ctx, cancel := CreateCtx().
		Set("User", "alice").
		Build()
	defer cancel()

	fmt.Println("User:", GetCtx(ctx, "User"))

	// Output:
	// User: alice
}

// This function is named ExampleCtxBuilder_Timeout()
// it with the Examples type.
func ExampleCtxBuilder_Timeout() {
	ctx, cancel := CreateCtx().
		Set("TraceID", "xyz").
		Timeout(2 * time.Second).
		Build()
	defer cancel()

	deadline, ok := ctx.Deadline()
	fmt.Println("Has deadline:", ok)
	fmt.Println("Timeout set to ~2s:", time.Until(deadline) <= 2*time.Second)

	// Output:
	// Has deadline: true
	// Timeout set to ~2s: true
}

// This function is named ExampleCtxBuilder_Build()
// it with the Examples type.
func ExampleCtxBuilder_Build() {
	builder := CreateCtx().
		Set("Key1", "Value1").
		Set("Key2", "Value2")

	ctx, cancel := builder.Build()
	defer cancel()

	fmt.Println("Key1:", GetCtx(ctx, "Key1"))
	fmt.Println("Key2:", GetCtx(ctx, "Key2"))

	// Output:
	// Key1: Value1
	// Key2: Value2
}
