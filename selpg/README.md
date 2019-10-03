# Selpg说明
## 实验报告

见CSDN博客，[传送门](https://blog.csdn.net/JKJKL1/article/details/101999617)

## Usage
> Program Name: cli
>
> Usage of cli:
>
>  -d, --destnation string   The destnatioon of printing
>
>  -e, --end int             The end page (default -1)
>
>  -f, --endOfPage           Defind the end symbol of a page
>
>  -l, --lineOfPage int      The length of a page (default 72)
>
>  -s, --start int           The start page (default -1)

## 文件说明

- **createInput.go** 用于创建初始的inputFile，也就是test.txt, 之后没有再用到，所以实际上没有什么用

- **selpg.go**  项目主要文件，包括flag的定义解析、文件读取输出

- **test.txt**  输入文件，为1~200的数字，每个数字一行，每10行有一个换页符'\f'，便于观察输出正确与否
