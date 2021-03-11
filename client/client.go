package client

import (
	"fmt"
	"os"
	"flag"
	"XianfengChain03/chain"
	"math/big"
)

/**
 * 客户端（命令行窗口工具），主要用户实现与用户进行动态交互
	① 将帮助信息等输出到控制台
	② 读取参数并解析，根据解析结果调用对应的项目功能
 */
type Client struct {
	Chain chain.BlockChain
}

/**
 * Client的run方法，是程序的主要处理逻辑入口
 */
func (client *Client) Run() {
	if len(os.Args) == 1 { //用户没有输入任何指令
		client.Help()
		return
	}
	//1、解析命令行参数
	command := os.Args[1]
	//2、确定用户输入的命令
	switch command {
	case CREATECHAIN:
		client.CreateChain()
	case GENERATEGENESIS:
		client.GenerateGensis()
	case ADDNEWBLOCK:
		client.AddNewBlock()
	case GETLASTBLOCK:
		client.GetLastBlock()
	case GETALLBLOCKS:
		client.GetAllBlocks()
	case GETBLOCKCOUNT:
		client.GetBlockCount()
	case HELP:
		client.Help()
	default:
		client.Default()
	}
	//3、根据不同的命令，调用blockChain的对应功能
	//4、根据调用的结果，将功能调用结果信息输出到控制台,提供给用户

}

func (client *Client) Default() {
	fmt.Println("go run main.go : Unknown subcommand.")
	fmt.Println("Use go run main.go help for more information.")
}

func (client *Client) GetBlockCount() {
	blocks, err := client.Chain.GetAllBlocks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("查询成功,当前共有%d个区块\n", len(blocks))
}

func (client *Client) GetAllBlocks() {
	if len(os.Args[2:]) > 0 {
		fmt.Println("抱歉，getallblocks不接收参数")
		return
	}

	allBlocks, err := client.Chain.GetAllBlocks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("成功获取到所有区块")
	for _, block := range allBlocks {
		fmt.Printf("区块%d,Hash:%x,数据:%s\n", block.Height, block.Hash, block.Data)
	}
}

func (client *Client) GetLastBlock() {
	set := os.Args[2:]
	if len(set) > 0 {
		fmt.Println("兄弟，你会错意了。getlastblock命令不是这么用的.")
		return
	}

	last := client.Chain.GetLastBlock()
	hashBig := new(big.Int)
	hashBig.SetBytes(last.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) > 0 {
		fmt.Println("查询到最新区块")
		fmt.Println("最新区块高度:", last.Height)
		fmt.Println("最新区块的内容:", string(last.Data))
		fmt.Printf("最新区块哈希:%x\n", last.Hash)
		fmt.Printf("前一个区块哈希:%x\n", last.PreHash)
		return
	}
	fmt.Println("抱歉，当前暂无最新区块")
	fmt.Println("请使用go run main.go generategensis生成创世区块")
}

func (client *Client) AddNewBlock() {
	addBlock := flag.NewFlagSet(ADDNEWBLOCK, flag.ExitOnError)
	data := addBlock.String("data", "", "区块存储的自定义内容")
	addBlock.Parse(os.Args[2:])

	//args := os.Args[2:]
	//1、从参数中取出所有以-开头的参数项
	//2、准备一个当前命令支持的所有参数的切片
	//系统预定义的： [-data]
	//用户输入的：  [-data -abcd]

	err := client.Chain.AddNewBlock([]byte(*data))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("恭喜，已成功创建新区块，并存储到文件中")
}
func (client *Client) GenerateGensis() {
	generateGensis := flag.NewFlagSet(GENERATEGENESIS, flag.ExitOnError)
	gensis := generateGensis.String("gensis", "", "创世区块中的自定义数据")
	generateGensis.Parse(os.Args[2:])
	//1、先判断是否已存在创世区块
	hashBig := new(big.Int)
	hashBig = hashBig.SetBytes(client.Chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 1 { //创世区块的hash值不为0，即有值
		fmt.Println("抱歉，创世区块已存在，无法覆盖写入")
		return
	}
	//2、如果创世区块不存在，才去调用creategenesis
	client.Chain.CreateGenesis([]byte(*gensis))
	fmt.Println("恭喜，创世区块创建并成功写入数据")
}

func (client *Client) CreateChain() {
	//

}

/**
 * 该方法用于向控制台输出项目的使用说明
 */
func (client *Client) Help() {
	fmt.Println("-------------Welcome to XianfengChain03 project-------------")
	fmt.Println()
	fmt.Println("USAGE：")
	fmt.Println("\tgo run main.go command [arguments]")
	fmt.Println()
	fmt.Println("AVAILABLE COMMANDS：")
	fmt.Println()
	fmt.Println("    " + CREATECHAIN + "       the command is used to create a new blockchain.")
	fmt.Println("    " + GENERATEGENESIS + "    generate a gensis block, use the gensis argument for the data.")
	fmt.Println("    addnewblock       create a new block, the argument is data.")
	fmt.Println()
	fmt.Println("Use go run main.go help for more information about a command.")
	fmt.Println()
}
