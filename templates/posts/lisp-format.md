**TLDR**: use fixed-size indentation ([2 spaces](/posts/spaces-tabs)) and place closing delimiters on their own lines.

## Disclaimer



This post comes from experience.


**self-inflicted problem** that **turns people away**

Not a fundamental problem with the `(inside_parens)` call style, just an impractical convention we don't have to follow.


## Problem

## Solution

```
(define (some-func . args)
        (let ((one 10)
              (two 20))
             (three four)
             (if five
                 (let ((six 60))
                   (seven eight)
                   (nine ten)))))
```

Racket actually allows the following:

```
(define (some-func . args)
  (define one 10)
  (define two 20)
  '(three four)
  (when 'five
    (define six 60)
    '(seven eight)
    '(nine ten)
  )
)
```


I opened the Racket source code, searched for `(let`, and one of the first results was this:

```scm
(define reverse-bytes
  (let ([pairs (let ([xs (bytes->list #"([{<")]
                     [ys (bytes->list #")]}>")])
                 (append (map cons xs ys) (map cons ys xs)))])
    (define (rev-byte b)
      (cond [(assq b pairs) => cdr]
            [else b]))
    (lambda (bs) (list->bytes (map rev-byte (reverse (bytes->list bs)))))))
```

```scm
(define reverse-bytes
  (block
    (define xs (bytes->list #"([{<"))
    (define ys (bytes->list #")]}>"))
    (define pairs (append (map cons xs ys) (map cons ys xs)))
    (define (rev-byte b)
      (cond
        (assq b pairs) => cdr
        else b
      )
    )
    (lambda (bs) (list->bytes (map rev-byte (reverse (bytes->list bs)))))
  )
)
```

```scm
(define (px . args)
  (let* ([args (let loop ([xs args])
                 (if (list? xs) (apply append (map loop xs)) (list xs)))]
         [args (map (lambda (x)
                      (cond [(bytes? x) x]
                            [(string? x) (string->bytes/utf-8 x)]
                            [(char? x) (regexp-quote (string->bytes/utf-8 (string x)))]
                            [(not x) #""]
                            [else (internal-error 'px)]))
                    args)])
    (byte-pregexp (apply bytes-append args))))
```
