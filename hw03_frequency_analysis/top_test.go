package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

var repetiveText = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut vulputate tellus
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut vulputate tellus
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut vulputate tellus
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut vulputate tellus
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut vulputate tellus
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut vulputate tellus
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut vulputate tellus
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut vulputate tellus
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut vulputate tellus`

func TestTop10(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "no words in empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:  "positive test",
			input: text,
			expected: []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			},
		},
		{
			name:  "separators",
			input: "i am\t\ta st\trange \n\t a\tnd st\tupid case",
			expected: []string{
				"a",
				"st",
				"am",
				"case",
				"i",
				"nd",
				"range",
				"upid",
			},
		},
		{
			name:  "equal_frequencies",
			input: repetiveText,
			expected: []string{
				"Etiam",
				"Lorem",
				"adipiscing",
				"amet,",
				"consectetur",
				"dolor",
				"elit.",
				"ipsum",
				"sit",
				"tellus",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, Top10(tc.input))
		})
	}
}
