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
		flag.NewFlagSet(CREATECHAIN, flag.ExitOnError)
	case GENERATEGENESIS:
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
	case ADDNEWBLOCK:
		fmt.Println("调用生成新区块的功能")
	case GETLATBLOCK:
		fmt.Println("获取最新区块的功能")
	case GETALLBLOCKS:
		fmt.Println("获取所有区块的功能")
	case GETBLOCKCOUNT:
		fmt.Println("获取区块的数量")
	case HELP:
		client.Help()
	default:
		fmt.Println("go run main.go : Unknown subcommand.")
		fmt.Println("Use go run main.go help for more information.")
	}
	//3、根据不同的命令，调用blockChain的对应功能
	//4、根据调用的结果，将功能调用结果信息输出到控制台,提供给用户

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
