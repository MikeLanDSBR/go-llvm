package llvm

import (
	"fmt"
	"strings"
)

func (t TypeKind) String() string {
	switch t {
	case VoidTypeKind:
		return "VoidTypeKind"
	case FloatTypeKind:
		return "FloatTypeKind"
	case DoubleTypeKind:
		return "DoubleTypeKind"
	case X86_FP80TypeKind:
		return "X86_FP80TypeKind"
	case FP128TypeKind:
		return "FP128TypeKind"
	case PPC_FP128TypeKind:
		return "PPC_FP128TypeKind"
	case LabelTypeKind:
		return "LabelTypeKind"
	case IntegerTypeKind:
		return "IntegerTypeKind"
	case FunctionTypeKind:
		return "FunctionTypeKind"
	case StructTypeKind:
		return "StructTypeKind"
	case ArrayTypeKind:
		return "ArrayTypeKind"
	case PointerTypeKind:
		return "PointerTypeKind"
	case VectorTypeKind:
		return "VectorTypeKind"
	case MetadataTypeKind:
		return "MetadataTypeKind"
	default: // Adicionado: Este caso captura qualquer TypeKind não mapeado
		return fmt.Sprintf("UNKNOWN_TYPE_KIND(%d)", t) // Retorna uma string descritiva em vez de pânico
	}
}

func (t Type) String() string {
	ts := typeStringer{s: make(map[Type]string)}
	return ts.typeString(t)
}

type typeStringer struct {
	s map[Type]string
}

func (ts *typeStringer) typeString(t Type) string {
	if s, ok := ts.s[t]; ok {
		return s
	}

	k := t.TypeKind()
	s := k.String()
	// Apenas remove "Kind" se a string não for "UNKNOWN_TYPE"
	if !strings.HasPrefix(s, "UNKNOWN_TYPE_KIND") {
		s = s[:len(s)-len("Kind")]
	}

	switch k {
	case ArrayTypeKind:
		s += fmt.Sprintf("(%v[%v])", ts.typeString(t.ElementType()), t.ArrayLength())
	case PointerTypeKind:
		// Para PointerTypeKind, verifica se ElementType() é válido antes de chamar typeString
		elementType := t.ElementType()
		if !elementType.IsNil() { // Verifica se o tipo do elemento não é nulo
			s += fmt.Sprintf("(%v)", ts.typeString(elementType))
		} else {
			s += "(nil)" // Ou outra representação para tipo de elemento nulo
		}
	case FunctionTypeKind:
		params := t.ParamTypes()
		s += "("
		if len(params) > 0 {
			s += fmt.Sprintf("%v", ts.typeString(params[0]))
			for i := 1; i < len(params); i++ {
				s += fmt.Sprintf(", %v", ts.typeString(params[i]))
			}
		}
		s += fmt.Sprintf("):%v", ts.typeString(t.ReturnType()))
	case StructTypeKind:
		if name := t.StructName(); name != "" {
			ts.s[t] = "%" + name
			s = fmt.Sprintf("%%%s: %s", name, s)
		}
		etypes := t.StructElementTypes()
		s += "("
		if n := len(etypes); n > 0 {
			s += ts.typeString(etypes[0])
			for i := 1; i < n; i++ {
				s += fmt.Sprintf(", %v", ts.typeString(etypes[i]))
			}
		}
		s += ")"
	case IntegerTypeKind:
		s += fmt.Sprintf("(%d bits)", t.IntTypeWidth())
	}

	ts.s[t] = s
	return s
}
