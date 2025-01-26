section .text
        global main

main:
  mov R12, 1
  mov R13, 2
  add RAX, R12, R13
  ret

        
;; hello world
;; section .text
;;     global main   ; Entry point for the linker (required by Windows)

;; main:
;;     mov rax, 2    ; Set the return value (exit code) in the RAX register
;;     ret           ; Return from the main function


;;         .global _main
;; .align 2

;; _main:
;;   mov R12, #1
;;   mov R13, #2
;;   add RAX, R12, R13
;;   ret

