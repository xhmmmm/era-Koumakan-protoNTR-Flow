﻿紅魔館protoNTR0.021対象一文字変数撲滅パッチ
===========================================

導入前の注意
------------
**このパッチはファイルの削除を要求するためただ上書きするだけでは動作しません。**
**必ずABLUP99.ERBを削除してください。**

**このパッチは一部不透明なライセンス問題を抱えている可能性があります。**

すでに適応している、またはこれから適応するパッチに一文字変数が含まれていないことを確認してください。
**このパッチは紅魔館protoNTR0.021が対象でそれ以外のパッチについては想定されていません。**
とくに、このパッチの後にパッチを当てるさいは注意してください。

このパッチを導入することによるゲーム性の変更は有りません。

概要
----
紅魔館protoNTR0.021から一文字変数を削除したうえで禁止します。
このパッチはコードを変更しただけなのでコメント部分に修正ミスがある可能性があります。
このパッチは他のライセンスに抵触しない範囲で自由に改変と配布が可能です。

手順
----
1. ABLUP99.ERBを削除してください。
2. このパッチのreadme以外のファイルを当ててください。

内容
----
このパッチは以下の関数を追加します。指定された用途では一文字変数ではなくこの関数を利用してください。
+ COMMON.ERB
    - `#FUNCTION @GET_STAINCOUNT(誰の部位, どこの部位 = -1, 誰が感じる)`
        + {誰の部位}の{どこの部位}に対し{誰が感じる}の感じる汚れ量を勝手に計算する関数
        + コマンド中に不潔量を計算するときはこの関数から汚れ量を取得してそれを基に計算してください

このパッチは以下の関数の返り値を変更します。
+ TRACHECK_ORGASM.ERB
    - `@ORGASM_ADD(奴隷, 調教者 = 0)`
        + 以前は常に1を返していたが計算の結果として取り出した快楽強度を返すように変更。
        + 浮気快楽強度を計算する際に使われる。

このパッチは以下の関数を削除します。指定された用途では代わりに使える関数を利用してください。
+ INFO.ERB
	- `@LIFE_BAR`
        + 未使用

このパッチは以下の変数を追加します。
+ TFLAG
    - 26=刻印取得時従順変化済みフラグ
        + コマンドを一回実行するごとにその回に刻印を取得した際一回だけ従順が上がることがあります。
		それが行われたことを保存します。
    - 160=実行値
        + 実行値計算の際初期化され実行値が入ります
    - 161=実行値表示用フラグ
        + 実行値計算の際それが一つ目の項かそうでないかを調べるために使われます

このパッチは以下の口上の中身が書き換えられています。
どれも一文字変数をLOCALに置き換えただけです。
**ライセンスが確認できなかったのでライセンスに抵触している可能性があります。**
問題が有ったらその変更を追認するか別の形に書き換えて再配布していただけたら幸いです。
+ 咲夜口上
+ チルノ口上
+ 大妖精口上
+ 魔理沙口上
