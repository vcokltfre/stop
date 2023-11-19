#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#define STACK_SIZE 256
#define CALL_STACK_SIZE 256

// #define DEBUG

#ifdef DEBUG
#define debug(...) printf(__VA_ARGS__)
#else
#define debug(...)
#endif

#define IHeaderHlt 0x00         // Halt
#define IHeaderDbg 0x01         // Debug
#define IHeaderMovLiteral 0x08  // Move value
#define IHeaderMovRegister 0x09 // Move register
#define IHeaderPush 0x10        // Push value
#define IHeaderDup 0x11         // Duplicate value
#define IHeaderDrop 0x12        // Drop value
#define IHeaderSwap 0x13        // Swap values
#define IHeaderLd 0x20          // Load value
#define IHeaderSt 0x21          // Store value
#define IHeaderAdd 0x30         // Add
#define IHeaderSub 0x31         // Subtract
#define IHeaderMul 0x32         // Multiply
#define IHeaderDiv 0x33         // Divide
#define IHeaderMod 0x34         // Modulo
#define IHeaderLabel 0xA0       // Label
#define IHeaderCall 0xA1        // Call
#define IHeaderJmp 0xA2         // Jump
#define IHeaderJmpZ 0xA3        // Jump if zero
#define IHeaderJmpNZ 0xA4       // Jump if not zero
#define IHeaderJmpP 0xA5        // Jump if positive
#define IHeaderJmpN 0xA6        // Jump if negative
#define IHeaderRet 0xA7         // Return
#define IHeaderPutN 0xB0        // Put number
#define IHeaderPutC 0xB1        // Put character

const unsigned int ISizeHlt = 1;         // {header}
const unsigned int ISizeDbg = 1;         // {header}
const unsigned int ISizeMovLiteral = 10; // {header, reg, value[8]}
const unsigned int ISizeMovRegister = 3; // {header, reg, source}
const unsigned int ISizePush = 9;        // {header, value[8]}
const unsigned int ISizeDup = 1;         // {header}
const unsigned int ISizeDrop = 1;        // {header}
const unsigned int ISizeSwap = 1;        // {header}
const unsigned int ISizeLd = 2;          // {header, reg}
const unsigned int ISizeSt = 2;          // {header, reg}
const unsigned int ISizeAdd = 1;         // {header}
const unsigned int ISizeSub = 1;         // {header}
const unsigned int ISizeMul = 1;         // {header}
const unsigned int ISizeDiv = 1;         // {header}
const unsigned int ISizeMod = 1;         // {header}
const unsigned int ISizeLabel = 3;       // {header, label[2]}
const unsigned int ISizeCall = 3;        // {header, label[2]}
const unsigned int ISizeJmp = 3;         // {header, label[2]}
const unsigned int ISizeJmpZ = 3;        // {header, label[2]}
const unsigned int ISizeJmpNZ = 3;       // {header, label[2]}
const unsigned int ISizeJmpP = 3;        // {header, label[2]}
const unsigned int ISizeJmpN = 3;        // {header, label[2]}
const unsigned int ISizeRet = 1;         // {header}
const unsigned int ISizePutN = 1;        // {header}
const unsigned int ISizePutC = 1;        // {header}

int64_t stack[STACK_SIZE];
uint64_t sp = 0;

uint64_t call_stack[CALL_STACK_SIZE];
uint64_t csp = 0;

uint64_t jumps[1 << 16];
uint64_t registers[16];

uint64_t ip = 0;

void build_jumps(uint8_t *buffer, long size)
{
    uint64_t ip = 0;

    while (ip < size)
    {
        switch (buffer[ip])
        {
        case IHeaderHlt:
            ip += ISizeHlt;
            break;
        case IHeaderDbg:
            ip += ISizeDbg;
            break;
        case IHeaderMovLiteral:
            ip += ISizeMovLiteral;
            break;
        case IHeaderMovRegister:
            ip += ISizeMovRegister;
            break;
        case IHeaderPush:
            ip += ISizePush;
            break;
        case IHeaderDup:
            ip += ISizeDup;
            break;
        case IHeaderDrop:
            ip += ISizeDrop;
            break;
        case IHeaderSwap:
            ip += ISizeSwap;
            break;
        case IHeaderLd:
            ip += ISizeLd;
            break;
        case IHeaderSt:
            ip += ISizeSt;
            break;
        case IHeaderAdd:
            ip += ISizeAdd;
            break;
        case IHeaderSub:
            ip += ISizeSub;
            break;
        case IHeaderMul:
            ip += ISizeMul;
            break;
        case IHeaderDiv:
            ip += ISizeDiv;
            break;
        case IHeaderMod:
            ip += ISizeMod;
            break;
        case IHeaderLabel:
            ip += ISizeLabel;
            jumps[(buffer[ip - 1] << 8) | buffer[ip - 2]] = ip;
            break;
        case IHeaderCall:
            ip += ISizeCall;
            break;
        case IHeaderJmp:
            ip += ISizeJmp;
            break;
        case IHeaderJmpZ:
            ip += ISizeJmpZ;
            break;
        case IHeaderJmpNZ:
            ip += ISizeJmpNZ;
            break;
        case IHeaderJmpP:
            ip += ISizeJmpP;
            break;
        case IHeaderJmpN:
            ip += ISizeJmpN;
            break;
        case IHeaderRet:
            ip += ISizeRet;
            break;
        case IHeaderPutN:
            ip += ISizePutN;
            break;
        case IHeaderPutC:
            ip += ISizePutC;
            break;
        default:
            printf("Error: unknown instruction %d\n", buffer[ip]);
            return;
        }
    }

    ip = 0;
}

void push(uint64_t value)
{
    if (sp >= STACK_SIZE)
    {
        printf("Error: stack overflow\n");
        exit(1);
    }

    stack[sp++] = value;
}

uint64_t pop()
{
    if (sp <= 0)
    {
        printf("Error: stack underflow\n");
        exit(1);
    }

    return stack[--sp];
}

uint64_t peek()
{
    if (sp <= 0)
    {
        printf("Error: stack underflow\n");
        exit(1);
    }

    return stack[sp - 1];
}

void call(uint64_t address)
{
    if (csp >= CALL_STACK_SIZE)
    {
        printf("Error: call stack overflow\n");
        exit(1);
    }

    call_stack[csp++] = address;
}

uint64_t ret()
{
    if (csp <= 0)
    {
        printf("Error: call stack underflow\n");
        exit(1);
    }

    return call_stack[--csp];
}

uint16_t read_u16(uint8_t *buffer, uint64_t offset)
{
    return (buffer[offset + 1] << 8) | buffer[offset];
}

int64_t read_i64(uint8_t *buffer, uint64_t offset) {
    return (
        (uint64_t)buffer[offset + 0] << 0 |
        (uint64_t)buffer[offset + 1] << 8 |
        (uint64_t)buffer[offset + 2] << 16 |
        (uint64_t)buffer[offset + 3] << 24 |
        (uint64_t)buffer[offset + 4] << 32 |
        (uint64_t)buffer[offset + 5] << 40 |
        (uint64_t)buffer[offset + 6] << 48 |
        (uint64_t)buffer[offset + 7] << 56
    );
}

void i_mov_literal(uint8_t *buffer)
{
    ip++;
    uint8_t reg = buffer[ip];
    ip++;
    registers[reg] = read_i64(buffer, ip);
    ip += 8;
}

void i_mov_register(uint8_t *buffer)
{
    ip++;
    uint8_t reg = buffer[ip++];
    uint8_t source = buffer[ip++];
    registers[reg] = registers[source];
}

void i_push(uint8_t *buffer)
{
    ip++;
    push(read_i64(buffer, ip));
    ip += 8;
}

void i_dup(uint8_t *buffer)
{
    ip++;
    push(peek());
}

void i_drop(uint8_t *buffer)
{
    ip++;
    pop();
}

void i_swap(uint8_t *buffer)
{
    ip++;
    uint64_t a = pop();
    uint64_t b = pop();
    push(a);
    push(b);
}

void i_ld(uint8_t *buffer)
{
    ip++;
    uint8_t reg = buffer[ip++];
    push(registers[reg]);
}

void i_st(uint8_t *buffer)
{
    ip++;
    uint8_t reg = buffer[ip++];
    registers[reg] = pop();
}

void i_add(uint8_t *buffer)
{
    uint64_t a = pop();
    uint64_t b = pop();
    push(a + b);

    ip += ISizeAdd;
}

void i_sub(uint8_t *buffer)
{
    uint64_t a = pop();
    uint64_t b = pop();
    push(a - b);

    ip += ISizeSub;
}

void i_mul(uint8_t *buffer)
{
    uint64_t a = pop();
    uint64_t b = pop();
    push(a * b);

    ip += ISizeMul;
}

void i_div(uint8_t *buffer)
{
    uint64_t a = pop();
    uint64_t b = pop();
    push(a / b);

    ip += ISizeDiv;
}

void i_mod(uint8_t *buffer)
{
    uint64_t a = pop();
    uint64_t b = pop();
    push(a % b);

    ip += ISizeMod;
}

void i_jmp(uint8_t *buffer)
{
    ip++;
    uint16_t addr = read_u16(buffer, ip);
    ip = jumps[addr];
}

void i_jmpz(uint8_t *buffer)
{
    ip++;
    uint16_t addr = read_u16(buffer, ip);
    ip += 2;
    if (pop() == 0)
    {
        ip = jumps[addr];
    }
}

void i_jmpnz(uint8_t *buffer)
{
    ip++;
    uint16_t addr = read_u16(buffer, ip);
    ip += 2;
    if (pop() != 0)
    {
        ip = jumps[addr];
    }
}

void i_jmpp(uint8_t *buffer)
{
    ip++;
    uint16_t addr = read_u16(buffer, ip);
    ip += 2;
    if (pop() > 0)
    {
        ip = jumps[addr];
    }
}

void i_jmpn(uint8_t *buffer)
{
    ip++;
    uint16_t addr = read_u16(buffer, ip);
    ip += 2;
    if (pop() < 0)
    {
        ip = jumps[addr];
    }
}

void i_call(uint8_t *buffer)
{
    ip++;
    uint16_t addr = read_u16(buffer, ip);
    call(ip+2);
    ip = jumps[addr];
}

void i_ret(uint8_t *buffer)
{
    ip = ret();
}

void i_putn(uint8_t *buffer)
{
    ip++;
    printf("%ld\n", pop());
}

void i_putc(uint8_t *buffer)
{
    ip++;
    printf("%c", (char)pop());
}

int run(uint8_t *buffer, long size)
{
    build_jumps(buffer, size);

    while (ip < size) {
        switch (buffer[ip]) {
        case IHeaderHlt:
            return pop();
        case IHeaderDbg:
            printf("ip=%ld, sp=%ld, csp=%ld\n", ip, sp, csp);
            ip += ISizeDbg;
            break;
        case IHeaderMovLiteral:
            debug("movl %d %ld\n", buffer[ip + 1], read_i64(buffer, ip + 2));
            i_mov_literal(buffer);
            break;
        case IHeaderMovRegister:
            debug("movr %d %d\n", buffer[ip + 1], buffer[ip + 2]);
            i_mov_register(buffer);
            break;
        case IHeaderPush:
            debug("push %ld\n", read_i64(buffer, ip + 1));
            i_push(buffer);
            break;
        case IHeaderDup:
            debug("dup\n");
            i_dup(buffer);
            break;
        case IHeaderDrop:
            debug("drop\n");
            i_drop(buffer);
            break;
        case IHeaderSwap:
            debug("swap\n");
            i_swap(buffer);
            break;
        case IHeaderLd:
            debug("ld %d\n", buffer[ip + 1]);
            i_ld(buffer);
            break;
        case IHeaderSt:
            debug("st %d\n", buffer[ip + 1]);
            i_st(buffer);
            break;
        case IHeaderAdd:
            debug("add\n");
            i_add(buffer);
            break;
        case IHeaderSub:
            debug("sub\n");
            i_sub(buffer);
            break;
        case IHeaderMul:
            debug("mul\n");
            i_mul(buffer);
            break;
        case IHeaderDiv:
            debug("div\n");
            i_div(buffer);
            break;
        case IHeaderMod:
            debug("mod\n");
            i_mod(buffer);
            break;
        case IHeaderLabel:
            debug("label %d\n", read_u16(buffer, ip + 1));
            ip += ISizeLabel;
            break;
        case IHeaderCall:
            debug("call %d\n", read_u16(buffer, ip + 1));
            i_call(buffer);
            break;
        case IHeaderJmp:
            debug("jmp %d\n", read_u16(buffer, ip + 1));
            i_jmp(buffer);
            break;
        case IHeaderJmpZ:
            debug("jmpz %d\n", read_u16(buffer, ip + 1));
            i_jmpz(buffer);
            break;
        case IHeaderJmpNZ:
            debug("jmpnz %d\n", read_u16(buffer, ip + 1));
            i_jmpnz(buffer);
            break;
        case IHeaderJmpP:
            debug("jmpp %d\n", read_u16(buffer, ip + 1));
            i_jmpp(buffer);
            break;
        case IHeaderJmpN:
            debug("jmpn %d\n", read_u16(buffer, ip + 1));
            i_jmpn(buffer);
            break;
        case IHeaderRet:
            debug("ret\n");
            i_ret(buffer);
            break;
        case IHeaderPutN:
            debug("putn\n");
            i_putn(buffer);
            break;
        case IHeaderPutC:
            debug("putc\n");
            i_putc(buffer);
            break;
        default:
            printf("Error: unknown instruction %d\n", buffer[ip]);
            return 1;
        }
    }

    return 0;
}

int main(int argc, char *argv[])
{
    if (argc < 2)
    {
        printf("Usage: %s <file>\n", argv[0]);
        return 1;
    }

    FILE *fp = fopen(argv[1], "rb");
    if (!fp)
    {
        printf("Error: could not open file %s\n", argv[1]);
        return 1;
    }

    fseek(fp, 0, SEEK_END);
    long size = ftell(fp);
    fseek(fp, 0, SEEK_SET);

    char *buffer = malloc(size);
    fread(buffer, size, 1, fp);
    fclose(fp);

    return run(buffer, size);
}
