package llvm_test

import (
	"testing"

	"github.com/taquion-lang/go-llvm" // importa o pacote principal
)

func TestPasses(t *testing.T) {
	llvm.InitializeNativeTarget()
	llvm.InitializeNativeAsmPrinter()

	ctx := llvm.NewContext()

	mod := ctx.NewModule("fac_module")

	fac_args := []llvm.Type{ctx.Int32Type()}
	fac_type := llvm.FunctionType(ctx.Int32Type(), fac_args, false)
	fac := llvm.AddFunction(mod, "fac", fac_type)
	fac.SetFunctionCallConv(llvm.CCallConv)
	n := fac.Param(0)

	entry := llvm.AddBasicBlock(fac, "entry")
	iftrue := llvm.AddBasicBlock(fac, "iftrue")
	iffalse := llvm.AddBasicBlock(fac, "iffalse")
	end := llvm.AddBasicBlock(fac, "end")

	builder := ctx.NewBuilder()
	defer builder.Dispose()

	builder.SetInsertPointAtEnd(entry)
	If := builder.CreateICmp(llvm.IntEQ, n, llvm.ConstInt(ctx.Int32Type(), 0, false), "cmptmp")
	builder.CreateCondBr(If, iftrue, iffalse)

	builder.SetInsertPointAtEnd(iftrue)
	res_iftrue := llvm.ConstInt(ctx.Int32Type(), 1, false)
	builder.CreateBr(end)

	builder.SetInsertPointAtEnd(iffalse)
	n_minus := builder.CreateSub(n, llvm.ConstInt(ctx.Int32Type(), 1, false), "subtmp")
	call_fac_args := []llvm.Value{n_minus}
	call_fac := builder.CreateCall(fac_type, fac, call_fac_args, "calltmp")
	res_iffalse := builder.CreateMul(n, call_fac, "multmp")
	builder.CreateBr(end)

	builder.SetInsertPointAtEnd(end)
	res := builder.CreatePHI(ctx.Int32Type(), "result")
	phi_vals := []llvm.Value{res_iftrue, res_iffalse}
	phi_blocks := []llvm.BasicBlock{iftrue, iffalse}
	res.AddIncoming(phi_vals, phi_blocks)
	builder.CreateRet(res)

	err := llvm.VerifyModule(mod, llvm.ReturnStatusAction)
	if err != nil {
		t.Errorf("Error verifying module: %s", err)
		return
	}

	targ, err := llvm.GetTargetFromTriple(llvm.DefaultTargetTriple())
	if err != nil {
		t.Error(err)
	}

	mt := targ.CreateTargetMachine(llvm.DefaultTargetTriple(), "", "", llvm.CodeGenLevelDefault, llvm.RelocDefault, llvm.CodeModelDefault)

	pbo := llvm.NewPassBuilderOptions()
	defer pbo.Dispose()

	t.Run("no error running default pass", func(t *testing.T) {
		err := mod.RunPasses("default<Os>", mt, pbo)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("errors on unknown pass name", func(t *testing.T) {
		err := mod.RunPasses("badpassnamedoesnotexist", mt, pbo)
		if err == nil {
			t.Error("expecting error but got none")
		}

		if err.Error() != "unknown pass name 'badpassnamedoesnotexist'" {
			t.Errorf("expected error about unknow pass name, instead got %s", err)
		}
	})
}
