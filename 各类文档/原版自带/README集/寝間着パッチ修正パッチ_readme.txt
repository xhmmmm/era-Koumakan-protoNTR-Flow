﻿-------------------------------------------------------------------------------------------

-era紅魔館protoNTR用寝間着パッチ修正パッチ-


簡単なパッチの説明
	着替えるパッチの修正パッチ
	

導入手順
	"era紅魔館protoNTR Ver.0.007"に"era紅魔館protoNTR修正パッチ131021_2"を適用した後に
	このパッチを入れてください
	
主な修正・追加処理
	・主人公が寝間着服装から変更されないのを修正
	・服装変更時に貞操帯が外れるのを修正
	・貧乳以下は寝間着時に上半身下着(以下ブラ)を着用しない
		但し、レミリアのような貧乳且つ私服でブラを着用している場合は寝間着時もブラを着用する
		
	・服装変更条件変更
		 起床後移動すると着替える
		これに加え
		 起床後部屋を出ず1時間経つと衣服変更(但し押し倒されている、押し倒している場合は除く)
	・私服時に靴を履かないキャラは寝間着のまま移動する際も靴を履かない
	・名称"貞操帯"が表示されるように
		
		
このパッチで使用している追加フラグ
	CFLAG:360 寝間着
	CFLAG:361 寝間着着用
	CFLAG:362 起床時間外
	CFLAG:363 起床後寝間着時間
	
ファイルの変更・追加・削除一覧
	BEFORETRAIN.ERB
		82行
			MASTER寝間着変更
			
	CLOTHES.ERB
		778行目
			"貞操帯"をもう一つ追加
		1267行 @CLOTHES_Preset_NIGHTWEAR(着用者)
			処理変更(貞操帯やブラ云々の不具合修正)
		1382行 @CLOTHES_Preset_NIGHTWEAR_S(着用者)
			寝間着が全裸なら移動時も靴を履かない
		1400行 @寝間着名称(ARG)
			一部文字変更
			
	MOVEMENT.ERB
		85行辺り
			起床後寝間着時間フラグの減少処理
			起床後部屋を出ず1時間経つと衣服変更
			
		1075行
			起床後寝間着時間フラグの設置
			
13/10/20 UP
13/10/24 修正パッチUP
	
	
	改造はご自由に
-------------------------------------------------------------------------------------------












男主人子が眼の前にいても気にせずに着替えるよ

