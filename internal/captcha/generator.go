package captcha

import (
	"fmt"
	"strconv"
	"strings"

	baseutil "github.com/imajinyun/go-knifer/internal/base"
)

// CodeGenerator 对应 hutool-captcha 中的 CodeGenerator 接口。
//
//	Generate 返回验证码原始串（写入图片）。
//	Verify   校验用户输入是否与原始串匹配；不同实现可有不同语义
//	         （RandomGenerator 直接比较，MathGenerator 需要对算式求值）。
type CodeGenerator interface {
	Generate() string
	Verify(code, userInput string) bool
}

// ---------------------------------------------------------------------------
// RandomGenerator 对应 hutool RandomGenerator
// ---------------------------------------------------------------------------

// RandomGenerator 随机字符验证码生成器。
type RandomGenerator struct {
	// BaseStr 字符集合，默认数字 + 大小写字母。
	BaseStr string
	// Length 验证码长度。
	Length int
}

// NewRandomGenerator 使用默认字符集（数字 + 大小写字母）创建生成器。
func NewRandomGenerator(length int) *RandomGenerator {
	return &RandomGenerator{BaseStr: baseutil.BaseCharNumberUC, Length: length}
}

// NewRandomGeneratorWithBase 自定义字符集与长度。
func NewRandomGeneratorWithBase(base string, length int) *RandomGenerator {
	return &RandomGenerator{BaseStr: base, Length: length}
}

// Generate 生成随机字符串。
func (g *RandomGenerator) Generate() string {
	base := g.BaseStr
	if base == "" {
		base = baseutil.BaseCharNumberUC
	}
	n := g.Length
	if n <= 0 {
		n = 4
	}
	runes := []rune(base)
	out := make([]rune, n)
	for i := 0; i < n; i++ {
		out[i] = runes[baseutil.RandomInt(len(runes))]
	}
	return string(out)
}

// Verify 大小写不敏感比较；输入空白返回 false（与 hutool 保持一致）。
func (g *RandomGenerator) Verify(code, userInput string) bool {
	if strings.TrimSpace(userInput) == "" {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(code), strings.TrimSpace(userInput))
}

// ---------------------------------------------------------------------------
// MathGenerator 对应 hutool MathGenerator
// ---------------------------------------------------------------------------

const mathOperators = "+-*"

// MathGenerator 算式验证码生成器，生成形如 "12+3 =" 的字符串，
// 校验时对其求值与用户输入对比。
type MathGenerator struct {
	// NumberLength 参与运算的数字最大位数（hutool 默认 2）。
	NumberLength int
	// ResultHasNegativeNumber 计算结果是否允许为负数（hutool 默认 true）。
	ResultHasNegativeNumber bool
}

// NewMathGenerator 默认: numberLength=2, 允许负数结果。
func NewMathGenerator() *MathGenerator {
	return &MathGenerator{NumberLength: 2, ResultHasNegativeNumber: true}
}

// NewMathGeneratorWith 自定义参数。
func NewMathGeneratorWith(numberLength int, resultHasNegativeNumber bool) *MathGenerator {
	if numberLength <= 0 {
		numberLength = 2
	}
	return &MathGenerator{NumberLength: numberLength, ResultHasNegativeNumber: resultHasNegativeNumber}
}

// Length 返回验证码渲染长度（与 hutool 一致：numberLength*2 + 2）。
func (g *MathGenerator) Length() int { return g.NumberLength*2 + 2 }

// Generate 生成 "a op b=" 格式算式（用空格右补齐，模拟 hutool padAfter 行为）。
func (g *MathGenerator) Generate() string {
	limit := g.limit()
	op := mathOperators[baseutil.RandomInt(len(mathOperators))]
	a := baseutil.RandomInt(limit)
	var b int
	if !g.ResultHasNegativeNumber && op == '-' {
		if a == 0 {
			b = 0
		} else {
			b = baseutil.RandomInt(a)
		}
	} else {
		b = baseutil.RandomInt(limit)
	}
	n1 := padRight(strconv.Itoa(a), g.NumberLength, ' ')
	n2 := padRight(strconv.Itoa(b), g.NumberLength, ' ')
	return fmt.Sprintf("%s%c%s=", n1, op, n2)
}

// Verify 对 code 求值并与用户输入比较。
func (g *MathGenerator) Verify(code, userInput string) bool {
	got, err := strconv.Atoi(strings.TrimSpace(userInput))
	if err != nil {
		return false
	}
	v, ok := evalMathExpr(code)
	if !ok {
		return false
	}
	return v == got
}

// limit 返回操作数上限：1 后跟 numberLength 个 0。
func (g *MathGenerator) limit() int {
	limit := 1
	for i := 0; i < g.NumberLength; i++ {
		limit *= 10
	}
	return limit
}

// padRight 将 s 用 c 右侧补齐到 n。
func padRight(s string, n int, c byte) string {
	if len(s) >= n {
		return s
	}
	pad := make([]byte, n-len(s))
	for i := range pad {
		pad[i] = c
	}
	return s + string(pad)
}

// evalMathExpr 解析 "a op b=" 形式（含空格补齐）的简单整数算式。
func evalMathExpr(s string) (int, bool) {
	s = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "="))
	for _, op := range []byte{'+', '-', '*'} {
		// 找到第一个非首字符位置上的运算符
		if i := strings.IndexByte(s, op); i > 0 {
			left := strings.TrimSpace(s[:i])
			right := strings.TrimSpace(s[i+1:])
			a, errA := strconv.Atoi(left)
			b, errB := strconv.Atoi(right)
			if errA != nil || errB != nil {
				return 0, false
			}
			switch op {
			case '+':
				return a + b, true
			case '-':
				return a - b, true
			case '*':
				return a * b, true
			}
		}
	}
	return 0, false
}
