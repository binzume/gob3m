# gob3m

近藤科学の [B3Mシリーズ](http://kondo-robot.com/product-category/servomotor/b3m) のサーボモーターをGolangから扱うライブラリ(の予定)．

まだ書きかけです．インターフェイスも仮です．

とりあえず動くもの：

- b3m.GetVersion モデル名・バージョン番号取得
- b3m.GetMode モード取得
- b3m.SetMode モード設定
- b3m.Reset サーボリセット
- b3m.SetPosition 目標位置設定
- b3m.GetCurrentPosition 現在位置取得


