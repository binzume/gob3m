# gob3m

近藤科学の [B3Mシリーズ](http://kondo-robot.com/product-category/servomotor/b3m) のサーボモーターをGolangから扱うライブラリ．

とりあえず動くもの：

- GetVersion モデル名・バージョン番号取得
- Reset サーボリセット
- GetMode モード取得
- SetMode モード設定
- SetPosition 目標位置設定
- SetVelocity 目標速度設定
- SetTorque 目標トルク設定
- GetCurrentPosition 現在位置取得
- 関数が用意されてないものは ReadMem/WriteMem で直接メモリを読み書きしてください

## Usage

examples/misc.go 参照．

```go
package main

import (
	"github.com/binzume/gob3m/b3m"
	"github.com/tarm/serial"
	"log"
)

func main() {
	s, err := serial.OpenPort(&serial.Config{Name: "COM3", Baud: 1500000})
	if err != nil {
		log.Fatal(err)
	}
	var id byte = 0

	servo := b3m.GetServo(s, id)

	err = servo.SetMode(b3m.ControlPosition | b3m.RunNormal)
	if err != nil {
		log.Fatal(err)
	}

	err = servo.SetPosition(500)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ok")
}
```
## TODO

- マルチモード対応
- タイムアウト処理まともにする
- パッケージ名変えるかも

## License

Copyright 2016 Kousuke Kawahira

Released under the MIT license

