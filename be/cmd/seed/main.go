package main

import (
	"fmt"
	"log"
	"math/rand"

	"server/config"
	"server/internal/model"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("config init failed: %v", err)
	}

	if err := config.InitDB(); err != nil {
		log.Fatalf("database init failed: %v", err)
	}
	defer config.CloseDB()

	db := config.DB()

	// 创建 5 个测试用户
	users := []model.User{
		{Username: "小红薯", Avatar: "https://i.pravatar.cc/150?img=1", Bio: "爱生活爱分享"},
		{Username: "美食达人", Avatar: "https://i.pravatar.cc/150?img=2", Bio: "吃遍天下无敌手"},
		{Username: "旅行博主", Avatar: "https://i.pravatar.cc/150?img=3", Bio: "一直在路上"},
		{Username: "时尚教主", Avatar: "https://i.pravatar.cc/150?img=5", Bio: "穿搭分享"},
		{Username: "健身教练", Avatar: "https://i.pravatar.cc/150?img=6", Bio: "带你练出好身材"},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			log.Printf("创建用户失败: %v", err)
		}
	}
	log.Printf("创建了 %d 个用户", len(users))

	// 示例图片 URL（来自 picsum.photos，随机不同图片）
	sampleImages := []string{
		"https://picsum.photos/800/600?random=1",
		"https://picsum.photos/800/600?random=2",
		"https://picsum.photos/800/600?random=3",
		"https://picsum.photos/800/600?random=4",
		"https://picsum.photos/800/600?random=5",
		"https://picsum.photos/800/600?random=6",
		"https://picsum.photos/800/600?random=7",
		"https://picsum.photos/800/600?random=8",
		"https://picsum.photos/800/600?random=9",
		"https://picsum.photos/800/600?random=10",
	}

	titles := []string{
		"今天的早餐太美味了！", "周末去爬山，呼吸新鲜空气",
		"分享一套超有效的减脂训练", "这件衣服太好看了，必须推荐",
		"旅行日记第一站：厦门", "家常菜也能做出高级感",
		"素颜妆教程，新手也能学会", "这本书真的强烈推荐",
		"露营装备清单分享", "如何养成早起习惯",
		"护肤品的正确涂抹顺序", "低成本改造出租屋",
		"烘焙小白第一次做蛋糕", "城市周边一日游推荐",
		"职场穿搭分享，干练又气质", "极简主义生活的好处",
		"春日穿搭灵感", "低成本养猫攻略",
		"周末早餐合集", "健身餐一周食谱",
	}

	bodies := []string{
		"今天尝试了一下新的做法，果然没有让我失望！强烈推荐给大家。",
		"周末终于有时间出去走走，天气也特别给力，心情超好。",
		"这套动作每天做一遍，坚持一个月效果明显，亲测有效。",
		"最近发现的宝藏单品，性价比超高，链接已经放在评论区啦。",
		"第一站选择了海边的小镇，真的太美了，随手一拍都是大片。",
		"简单几步，在家也能做出餐厅级别的味道，全家都爱吃。",
		"新手友好，手把手教你画出自然又好看的妆容。",
		"读完之后受益匪浅，重新审视了自己的生活方式。",
		"清单整理好了，打算去露营的朋友们可以直接抄作业。",
		"坚持早起一段时间后，整个人精神状态都变好了，给大家分享我的方法。",
		"很多人顺序都用错了，正确顺序能让效果翻倍。",
		"改造只花了不到500块，效果超出预期！",
		"第一次做没想到这么成功，配方简单，味道超好。",
		"距离市区开车只要1小时，适合周末亲子游。",
		"职场穿搭不需要花大钱，基础款也能穿出高级感。",
		"极简生活后减少了80%的物品，生活反而更轻松了。",
		"春天就要穿明亮的颜色，这几套搭配太爱了。",
		"养猫成本其实可以很低，这些平替一定要知道。",
		"周末早起的动力就是给自己做一顿丰盛的早餐。",
		"一周健身餐计划，每天不重样，营养又健康。",
	}

	// 批量创建 200 条笔记
	noteCount := 200
	for i := 0; i < noteCount; i++ {
		user := users[rand.Intn(len(users))]
		title := titles[rand.Intn(len(titles))]
		body := bodies[rand.Intn(len(bodies))]

		// 随机选 1-5 张图片
		imageCount := rand.Intn(5) + 1
		images := make([]model.NoteImage, imageCount)
		for j := 0; j < imageCount; j++ {
			images[j] = model.NoteImage{
				URL:    sampleImages[rand.Intn(len(sampleImages))],
				Width:  800,
				Height: 600,
			}
		}

		note := model.Note{
			Title:        fmt.Sprintf("%s #%d", title, i+1),
			Body:         body,
			UserID:       user.ID,
			LikeCount:    uint(rand.Intn(1000)),
			CollectCount: uint(rand.Intn(500)),
			Images:       images,
		}

		if err := db.Create(&note).Error; err != nil {
			log.Printf("创建笔记失败: %v", err)
			continue
		}

		if (i+1)%20 == 0 {
			log.Printf("已创建 %d/%d 条笔记", i+1, noteCount)
		}
	}

	log.Printf("完成！共创建了 %d 条笔记", noteCount)
}
