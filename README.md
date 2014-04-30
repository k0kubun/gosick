Dear All

Here is the description of programming project for newcomers.
(written in Japanese, sorry)

新4年生むけに毎年行っているプログラミング課題の案内です．
言語処理系・関数プログラミング・記号処理・システムプログラミングに慣れ親しむのが目的です．
M1以上の人が参加してもかまいません．


コースA: Schemeで実装する
* 以下に述べるSchemeサブセットの機能全部を実装すること
* call/cc, 限定継続, unwind-protect のいずれかを実装すること
* 実装したSchemeサブセット自身で，実装した処理系を動作させることができるようにすること

コースB: Scheme以外の言語で実装する（S式の構文解析も含める）
* TCOとマクロはオプション
* CLOS風のオブジェクト指向言語機構, アクターモデル風の並行計算機構のいずれかを実装する

コースC: C/C++で実装し，mbed LPC1768で動作させる．
* TCOとマクロはオプション
* GCを実装すること
* ホスト側である程度前処理して（例えば仮想機械語にコンパイルして），その出力をmbedで動作させてもよい
その場合ホスト側はScheme等で実装してもかまわない．
* mbedアプリケーションボード，あるいは他のセンサなどを使えるようにしてみる（必要な部品や工具は用意します）


いずれのコースとも，以下の要件をみたすこと．
・構文に合致しない入力はエラーとすること．
・実行時エラーによってインタプリタが停止してしまわないようにすること(できる範囲で)．


以上です．期間約1ヶ月程度とします．完成したものはゼミの時間にデモをしてもらいます．
わからないことは遠慮なく私や先輩にきいてください．


実装するSchemeサブセットについて

A. 構文

非終端記号:
 大文字[A-Z]で始まる名前

終端記号:
 開き括弧 (
 閉じ括弧 )
 ピリオド .
 小文字[a-z]で構成される名前
 シングルクオート ' で囲まれた文字列

X*  要素Xの0個以上のくりかえし
X+  要素Xの1個以上のくりかえし
[X] 要素Xは高々1個

Toplevel ::= Exp
       | Define
       | (load String)                      ファイルの内容を読み込む

Define ::= (define Id Exp)
     | (define (Id Id* [. Id]) Body)

Exp ::= Const                                   定数
  | Id                                      変数
  | (lambda Arg Body)                       λ抽象
  | (Exp Exp*)                              関数適用
  | (quote S-Exp)                           クオート (注1)
  | (set! Id Exp)                           代入
  | (let [Id] Bindings Body)                let
  | (let* Bindings Body)                    let* (注4)
  | (letrec Bindings Body)                  letrec
  | (if Exp Exp [Exp])                      条件式(if)
  | (cond (Exp Exp+)* [(else Exp+)])        条件式(cond) (注5)
  | (and Exp*)                              条件式(and)
  | (or Exp*)                               条件式(or)
  | (begin Exp*)                            逐次実行
  | (do ((Id Exp Exp)*) (Exp Exp*) Body)    繰り返し

Body ::= Define* Exp+

Arg ::= Id
 | (Id* [Id . Id])

Bindings ::= ((Id Exp)*)

S-Exp ::= Const
   | Id
   | (S-Exp* [S-Exp . S-Exp])

Const ::= Num
   | Bool
   | String
   | ()

Num: 最低限10進整数はサポートすること

Bool ::= '#t'
  | '#f'

String: ダブルクオートで囲まれた文字列 (注1)

Id ::= [0-9A-Za-z!$%&*+-./<=>?@^_]+   (注3)

(注1) (quote X) だけでなく 'X という記法もサポートすること．
(注2) バックスラッシュ \ によるエスケープをサポートすること．
(注3) 数値となるものをのぞくこと．
(注4) let* は let の0個以上の繰り返しではなく let* という名前である．
(注5) 節は(else節を含めて)最低1個はあること．else節は高々1個とすること．

B. 関数（最低限）

整数
number?, +, -, *, /, =, <, <=, >, >=

リスト
null?, pair?, list?, symbol?,
car, cdr, cons, list, length, memq, last, append,
set-car!, set-cdr!

ブール値
boolean?, not

文字列
string?, string-append,
symbol->string, string->symbol, string->number, number->string

関数
procedure?

比較
eq?, neq?, equal?

その他
load

C. 末尾呼び出しの最適化
末尾呼び出しの最適化(Tail Call Optimization, TCO)とは，関数本体の
末尾の位置(tail position)に関数呼び出しがある場合，コールスタックを
消費せずにその呼び出された関数に制御を移すことです．例えば

(define (even? x)
(if (= x 0) #t (odd? (- x 1))))
(define (odd? x)
(if (= x 1) #t (even? (- x 1))))

のような関数 even?, odd? の場合，0以上の引数に対しスタックオーバー
フローを起こさずに計算できなくてはなりません．

D. マクロ
マクロは以下のように定義します．

(define-macro (positive x)
(list '> x 0))

プログラム中に (positive e) という式が出てきた場合，これが
(list '> e 0) の実行結果で置き換えられ（これを「マクロが展開される」
といいます），その後で式の評価が行われます．例えば
(if (positive (+ (g x) 1)) (foo) (bar))
という式の場合，まずマクロが展開されて
(if (> (+ (g x) 1) 0) (foo) (bar))
という形の式が得られて，それからこの式が評価されます．

次に，以下のような関数定義との違いについて説明しておきます．

(define (positive x)
(> x 0))

関数定義の場合は，(positive e) という式を評価する場合は，まず
eを評価して，その結果の値が関数 positive の引数になります．一方
もしこれがマクロなら，e は評価されずに (> e 0) という形の式に
展開され，それから展開後の式が評価されます．

引数を評価しなくてすむので，関数定義と違って以下のような定義も
可能です．

(define (let*-expander vars body)
(if (null? vars)
 (cons 'begin body)
 (list 'let (list (car vars)) (let*-expander (cdr vars) body))))

(define-macro (let* vars . body)
(let*-expander vars body))

これは let を使って let* を定義した例です．このようにマクロを
使うことで，さまざまな構文をユーザが自分で定義することができます．

Schemeには define-syntax (let-syntax) というマクロ定義機構がありますが，
これは展開の前後で変数名がぶつからないようにする hygienic マクロと呼ばれる
もので，ちょっと実装は難しいと思います．今回実装してもらうのは Common-Lisp
風のマクロです．ちなみにCommon-Lispだと上の例は以下のようになります．

(defun let*-expander (vars body)
(if (null vars)
 (cons 'begin body)
 (list 'let (list (car vars)) (let*-expander (cdr vars) body))))

(defmacro let* (vars &rest body)
(let*-expander vars body))

以上．
