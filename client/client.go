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
	case SENDTRASACTION: //发送一笔新交易
		client.SendTransaction()
	case GETBALANCE: //获取地址的余额功能
		client.GetBalance()
	case GETLASTBLOCK:
		client.GetLastBlock()
	case GETALLBLOCKS: //获取所有的区块信息并打印输出给用户
		client.GetAllBlocks()
	case GETBLOCKCOUNT:
		client.GetBlockCount()
	case GETNEWADDRESS:
		client.GetNewAddress()
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

/**
 * 该方法用于生成一个新的地址
 * 地址的生成规则参考比特币地址的生成步骤
 */
func (client *Client) GetNewAddress() {
	getNewAddress := flag.NewFlagSet(GETNEWADDRESS, flag.ExitOnError)
	err := getNewAddress.Parse(os.Args[2:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//判断用户是否输了内容，getnewaddress命令不需要输入参数
	if len(os.Args[2:]) > 0 {
		fmt.Println("生成地址不需要参数，请重试！")
		return
	}
	address, err := client.Chain.GetNewAddress()
	if err != nil {
		fmt.Println("抱歉，地址生成错误，请重试。错误信息如下：", err.Error())
		return
	}
	fmt.Println("生成的地址是：", address)
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
		fmt.Printf("区块高度%d,区块hash：%x\n", block.Height, block.Hash)
		for index, tx := range block.Txs { //遍历区块中的每一笔交易
			fmt.Printf("区块%d的第%d笔交易,交易hash是：%x\n", block.Height, index, tx.TxHash)
			fmt.Println("\t该笔交易的交易输入：")
			for inputIndex, input := range tx.Inputs { //遍历交易的交易输入
				fmt.Printf("\t\t第%d个交易输入,花的钱来自%x中的第%d个输出\n", inputIndex, input.TxId, input.Vout)
			}
			fmt.Println("\t该笔交易的交易输出：")
			for outputIndex, output := range tx.Outputs { //遍历交易的交易输出
				fmt.Printf("\t\t第%d个交易输出，转给%s一笔面额为%f的钱\n", outputIndex, output.ScriptPub, output.Value)
			}
		}
		fmt.Println()
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
		fmt.Printf("最新区块哈希:%x\n", last.Hash)
		fmt.Printf("前一个区块哈希:%x\n", last.PreHash)
		fmt.Println("最新区块的交易:", last.Txs)
		return
	}
	fmt.Println("抱歉，当前暂无最新区块")
	fmt.Println("请使用go run main.go generategensis生成创世区块")
}

/**
 * 获取地址的余额的功能
 */
func (client *Client) GetBalance() {
	var address string
	getbalance := flag.NewFlagSet(GETBALANCE, flag.ExitOnError)
	getbalance.StringVar(&address, "address", "", "要查询的地址")
	getbalance.Parse(os.Args[2:])

	if len(address) == 0 {
		fmt.Println("请输入要查询的地址")
		return
	}
	totalbalance := client.Chain.GetBalance(address)
	fmt.Printf("地址%s的余额是%f\n ", address, totalbalance)
}

/**
 * 发送一笔新的交易
 */
func (client *Client) SendTransaction() {
	addBlock := flag.NewFlagSet(SENDTRASACTION, flag.ExitOnError)

	from := addBlock.String("from", "", "发起者地址")
	to := addBlock.String("to", "", "接收者地址")
	value := addBlock.String("value", "", "转账数量")
	addBlock.Parse(os.Args[2:])
	err := client.Chain.SendTransaction(*from, *to, *value)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("恭喜，已成功发送交易")
}

func (client *Client) GenerateGensis() {
	generateGensis := flag.NewFlagSet(GENERATEGENESIS, flag.ExitOnError)
	address := generateGensis.String("address", "", "用户指定的矿工地址")
	generateGensis.Parse(os.Args[2:])

	//1、先判断是否已存在创世区块
	hashBig := new(big.Int)
	hashBig = hashBig.SetBytes(client.Chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 1 { //创世区块的hash值不为0，即有值
		fmt.Println("抱歉，已有coinbase交易，暂不能重复构建")
		return
	}

	//2、如果创世区块不存在，才去调用creategenesis
	coinbaseHash, err := client.Chain.CreateCoinbase(*address)
	if err != nil {
		fmt.Println("抱歉，构建coinbase遇到错误，请重试。错误是：", err.Error())
		return
	}
	fmt.Printf("恭喜，COINBASE交易创建成功，交易hash是：%x\n", coinbaseHash)
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
	fmt.Println("    sendtransaction            create a new transaction, the argument is -from -to and -value.")
	fmt.Println("    " + GETNEWADDRESS + "    the command is used to generate a new address by bitcoin algorithm")
	fmt.Println()
	fmt.Println("Use go run main.go help for more information about a command.")
	fmt.Println()
}
