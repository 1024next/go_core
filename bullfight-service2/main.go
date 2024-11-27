package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Card struct {
	Suit  string
	Name  string
	Value int
}

type Player struct {
	Name     string
	Cards    []Card
	HandType string
	Multiple int
}

var suits = []string{"♠", "♥", "♦", "♣"}
var names = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
var values = map[string]int{
	"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "10": 10,
	"J": 10, "Q": 10, "K": 10, "A": 1,
}

var deck []Card
var players []Player

// 初始化牌堆
func initDeck() {
	for _, suit := range suits {
		for _, name := range names {
			deck = append(deck, Card{
				Suit:  suit,
				Name:  name,
				Value: values[name],
			})
		}
	}
}

// 洗牌
func shuffleDeck() {
	rand.Seed(time.Now().UnixNano())
	for i := len(deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

// 发牌
func dealCards() {
	for i := 0; i < 3; i++ {
		player := Player{
			Name: fmt.Sprintf("Player %d", i+1),
		}
		players = append(players, player)
	}

	dealer := Player{
		Name: "Dealer",
	}
	players = append(players, dealer)

	// 轮流发牌给玩家，每个玩家发5张牌
	for i := 0; i < 5; i++ {
		for j := 0; j < len(players); j++ {
			players[j].Cards = append(players[j].Cards, deck[i*len(players)+j])
		}
	}
}

// 判断是否是花牌
func isFlower(card Card) bool {
	return card.Name == "J" || card.Name == "Q" || card.Name == "K"
}

// 判断牌型
func checkHand(cards []Card) (string, int) {
	// 检查五小牛
	if allLessThanFive(cards) && sum(cards) <= 10 {
		return "五小牛", 10
	}
	// 检查四炸
	if isFourOfAKind(cards) {
		return "四炸", 8
	}
	// 检查五花牛
	if allFlowers(cards) {
		return "五花牛", 7
	}
	// 检查四花牛
	if isFourFlowers(cards) {
		return "四花牛", 6
	}
	// 检查牛牛
	if cow, isCow := checkCow(cards); isCow && cow == 0 {
		return "牛牛", 5
	}
	// 检查牛9
	if cow, isCow := checkCow(cards); isCow && cow == 9 {
		return "牛9", 4
	}
	// 检查牛8
	if cow, isCow := checkCow(cards); isCow && cow == 8 {
		return "牛8", 3
	}
	// 检查牛7
	if cow, isCow := checkCow(cards); isCow && cow == 7 {
		return "牛7", 2
	}
	if cow, isCow := checkCow(cards); isCow {
		return fmt.Sprintf("牛%d", cow), 1
	}

	// 没牛
	return "没牛", 1
}

// 检查五张牌的点数都小于5
func allLessThanFive(cards []Card) bool {
	for _, card := range cards {
		if card.Value >= 5 {
			return false
		}
	}
	return true
}

// 计算5张牌的总和
func sum(cards []Card) int {
	total := 0
	for _, card := range cards {
		total += card.Value
	}
	return total
}

func isFourOfAKind(cards []Card) bool {
	// 使用 map 来记录每种花牌的出现次数
	flowerCards := make(map[string]int)
	for _, card := range cards {
		if card.Name == "J" || card.Name == "Q" || card.Name == "K" {
			flowerCards[card.Name]++
		}
	}
	// 如果有四张 J、Q 或 K 中的某种花牌，则为四炸
	for _, count := range flowerCards {
		if count == 4 {
			return true
		}
	}
	return false
}

// 检查是否所有牌都是花牌（J, Q, K）
func allFlowers(cards []Card) bool {
	for _, card := range cards {
		if !isFlower(card) {
			return false
		}
	}
	return true
}

// 检查是否四张是花牌，剩下的1张是10
func isFourFlowers(cards []Card) bool {
	flowerCount := 0
	otherCount := 0
	for _, card := range cards {
		if isFlower(card) {
			flowerCount++
		} else if card.Value == 10 {
			otherCount++
		}
	}
	return flowerCount == 4 && otherCount == 1
}

// 判断是否为牛几
func checkCow(cards []Card) (int, bool) {
	// 遍历所有组合
	for i := 0; i < len(cards)-3+1; i++ {
		for j := i + 1; j < len(cards)-3+2; j++ {
			for m := j + 1; m < len(cards)-3+3; m++ {
				// 获取组合 [cards[i], cards[j], cards[m]]
				combination := []Card{cards[i], cards[j], cards[m]}
				// 计算组合的数值总和
				sum := combination[0].Value + combination[1].Value + combination[2].Value

				// 如果和是10的倍数，打印并返回
				if sum%10 == 0 {
					// 输出组合
					remaining := []Card{}
					// 找出剩余的卡片
					for _, card := range cards {
						if card != combination[0] && card != combination[1] && card != combination[2] {
							remaining = append(remaining, card)
						}
					}
					remainingSum := (remaining[0].Value + remaining[1].Value) % 10
					// 打印组合及剩余卡片
					fmt.Printf("组合: %v, 剩余: %v, 和: %d\n", combination, remaining, remainingSum)
					return remainingSum, true
				}
			}
		}
	}
	// 没有牛
	return 0, false
}

// 游戏逻辑
func gameLogic() {
	// 清空之前的游戏结果
	players = nil
	deck = nil

	// 初始化新的牌堆和玩家
	initDeck()
	shuffleDeck()
	dealCards()

	// 计算每个玩家的牌型及倍数
	for i := range players {
		handType, multiple := checkHand(players[i].Cards)
		players[i].HandType = handType
		players[i].Multiple = multiple
	}

	// 输出每个玩家的牌型和得分
	for _, player := range players {
		fmt.Printf("%s 的牌: ", player.Name)
		for _, card := range player.Cards {
			fmt.Printf("%s%s ", card.Suit, card.Name)
		}
		fmt.Printf("牌型: %s, 倍数: %d,", player.HandType, player.Multiple)
	}
}

// Gin 路由处理函数
func startGame(c *gin.Context) {
	gameLogic()
	// 返回游戏结果
	results := []string{}
	for _, player := range players {
		result := fmt.Sprintf("%s 的牌: ", player.Name)
		for _, card := range player.Cards {
			result += fmt.Sprintf("%s %s ", card.Suit, card.Name)
		}
		result += fmt.Sprintf("牌型: %s, 倍数: %d", player.HandType, player.Multiple)
		results = append(results, result)
	}
	c.JSON(http.StatusOK, gin.H{"game_results": results})
}

func main() {
	r := gin.Default()

	r.GET("/start", startGame)

	r.Run(":8080")
}
