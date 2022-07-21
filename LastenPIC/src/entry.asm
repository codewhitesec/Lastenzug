extern lastenzug
global entry

segment .text

entry:
    push rdi                    
    mov rdi, rsp               
    and rsp, byte -0x10       
    sub rsp, byte +0x20      
    call lastenzug              
    mov rsp, rdi           
    pop rdi               
    ret                  
