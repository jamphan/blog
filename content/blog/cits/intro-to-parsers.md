---
title: "Intro to Parsers in Compilers"
date: 2020-10-10T09:54:50+11:00
slug: "intro-to-parsers"
description: ""
keywords: ["grammars", "parsers", "compilers"]
draft: false
tags: ["compiler"]
math: false
toc: true
---

Python 3.9.0 was released recently (Oct 9th) and one of the changes was in [PEP 617 -- New PEG parser for CPython](https://www.python.org/dev/peps/pep-0617/).

This caught my interest, I wanted to know what a PEG parser is and why it was elected for the new 3.9.0 release over the LL-parser.

Naturally, this took me down a rabbit hole of reading; here are my notes on the relevant topics in understanding a PEG parser.

## The broader picture, where parsers sit in compilers.

{{< tldr "parsers are a component of the compiler 'front-end', helping to produce an intermediate representation of the code before the compiler synthesises the executable machine-code by performing syntax analysis." >}}

A compiler, in the context of computer systems, translates code written in one language (the source language) into another language (the target language). Invariably, compilers are used to translate human-readable source-code written in a particular programming language such as C into lower-level machine-language code, which can then be executed by the computer to perform the desired task(s).

Broadly speaking, a compiler undergoes two distinct phases, each performed by a corresponding category of sub-systems[^tu1]:

- The **Analysis phase** is performed by the **front-end** sub-systems, dividing the source code into its core parts and verifying for lexical, grammar and syntax errors. The key output of the front-end is an intermediate representation (IR) of the program which, along with a symbol table, is fed into the second phase to generate the actual program. Analysis consists of:

    - Lexical analysis
    - Syntax analysis (parser)
    - Semantic analysis
    - Intermediate code representation generation

- The **Synthesis phase** is performed by the **back-end** sub-systems. As the name suggests, this phase is what generates the code in the target language; in the case of machine-code, this generates the executable program. Synthesis consists of:

    - Machine independent code optimiser
    - Code generator
    - Machine dependent code optimiser

{{< figure src="/cits/grammars-intro_compilerArch.jpg" title="Compiler High-Level Architecture" attr="(from tutorialspoint.com)" attrlink="https://www.tutorialspoint.com/compiler_design/compiler_design_architecture.htm" >}}

In this high-level architectural view of the compiler, the parser is a front-end component. In particular, it is responsible for verifying the syntax of the input.

## A deeper diver into the compiler front-end

{{< tldr "syntax analysis ensures the expressions are syntactically correct; i.e. the arrangement of tokens is valid. However, it does not verify the expression make any sense." >}}

The front-end ensures there are no syntactical errors in the source-code, and simplifies retargeting (language translation) by providing a intermediate representation for the synthesis phase. The analysis within the front-end consists of:

1. **Lexical analysis (scanning)**: is the process of convertering the input of simple character sequences into a list of *tokens* of different kinds (such as numerical and string constants, variable identifiers, and programming language keywords)[^gr1]. For example, in C, the variable declaration line `int value = 100;` contains the tokens: `int` (keyword), `value` (identifier), `=` (operator), `100` (constant), and `;` (symbol)[^tu1].
2. **Syntax analysis (parsing)**: the lexer (or scanner) helps identifies tokens from the input source code, but it cannot check the syntax (the arrangement of words or phrases - i.e. tokens - to create well-formed sentences in a language) due to its limitations. The parser generates a syntax tree (or parse tree) and token arrangements are checked against the language's grammar; the parser checks the expressions made by the tokens are syntactically correct. For example, `10 = int m;` is syntactically invalid despite all characters matching some valid token.
3. **Semantic analysis**: While lexical analysis can identify tokens, and syntax analysis identify the correct arrangement of these tokens, the compiler must also analyse the underlying meaning of these constructs. Semantic analysis determines whether or not the syntax structure derive any meaning, for example, the statement `int a = "Hello, World!"` whilst syntactically correct, is semantically invalid as the type of assignment differs. Some more examples of what semantic analysis performs: scope resolution, type checking, array-bound checking[^tu1].

Through this dive into the compiler front-end, we have a better understanding of a parser and what it does - syntax analysis. We know what it takes in as input (tokens from the lexer), and what it generates as output (a parse tree). We know what it does (check syntax) and what it does not (check for meaning).

## What is Syntax? Context-Free Grammars and Productions

A compiler's parser has the primary responsiblity of recognising syntax; determining if the program being compiled contains valid sentences in the syntactic model of the programming language.
The syntactic model is expressed in a formal *Grammar*[^co2]. For some formal Grammar *G*, we say *G* dervies some string, if that string is syntactically valid in that language.

*Parsing* is the process of building a proof that the program being compiled, and all of its streams of words, can be derived by the programming language's formal Grammar.

- A **Grammar** is a set of rules for generating sentences in a language[^co1]; a formal method to describe a (textual) language.
- A **Parser** tests whether a text conforms to a grammar and if so, turns the valid text into a parse tree.

## Parse Trees

A parse tree is a graphical representation for the derivation, or parse, that corresponds to the input program[^sc1].

## References

<!-- References -->
[^co1]: (Lecture Notes, 2006) **Cornell University, Computer Science**. [Grammars and Parsing, Lecture 5](https://www.cs.cornell.edu/courses/cs211/2006sp/Lectures/L04-Parsing/L5cs211fa05.pdf), *CS-211: Computers and Programming*.

[^co2]: (Book, 2012). **K.D. Cooper and L. Torczon**. [Engineering a Compiler](https://www.sciencedirect.com/book/9780120884780/engineering-a-compiler), Second Edition, ISBN 978-0-12-088478-0.

[^gr1]: (Book, 1998) **R. Greenlaw and H.J. Hoover**. [Fundamentals of the Theory of Computation: Principles and Practice](https://www.sciencedirect.com/book/9781558605473/fundamentals-of-the-theory-of-computation-principles-and-practice#book-info), ISBN 978-1-55860-547-3.

[^od1]: (Lecture Notes, n.d.) **Old Dominion University, Dept. of Computer Science**. [Brief notes on parsing](https://www.cs.odu.edu/~toida/nerzic/390teched/cfl/Parsing/index.html), *CS-390: Introduction to Theoretical Computer Science/Theory of Computation*.

[^sc1]: (Webpage, n.d.) **Science Direct**. [Parse Tree](https://www.sciencedirect.com/topics/computer-science/parse-tree).

[^tu1]: (Webpage, n.d.) **Tutorialspoint.com**. [Compiler Design Tutorial](https://www.tutorialspoint.com/compiler_design/)

[^uf1]: (Lecture Notes, n.d.): http://people.cs.vt.edu/prsardar/classes/cs3304-Spr19/lectures/CS3304-9-LanguageSyntax-2.pdf
