# erc20_score项目说明
# 1、项目需求：

    1、部署一个带mint和burn功能的erc20合约，铸造销毁几个token，转移几个token，来构造事件
    2、使用go语言写一个后端服务来追踪合约事件，重建用户的余额
    3、以太坊延迟六个区块，确保区块链不会回滚
    4、加上积分计算功能，起一个定时任务，每小时根据用户的余额来计算用户的积分，暂定积分是余额*0.05
    5、要记录用户的所有余额变化，根据这个来计算积分，这样更准确一些
    6、需要维护一下用户的总余额表以及总积分表，还有一个用户的余额变动记录表
    7、需要支持多链逻辑，比如支持sepolia， base sepolia
	8、考虑一个场景，如果程序错误了，或者rpc有问题，导致好几天没有计算积分。此时应该如何正确回溯？


```c
    举个例子：
    用户在15：00的时候0个token，15：10分有100个token，15：30有200个token
    计算积分的时候，需要考虑用户的余额变化
    
    比如此时是16：00启动定时任务了来计算积分，应该是100*0.05*20/60+200*0.05*30/60
    
```


# 2、项目结构和说明
## 2.1 erc20_score_contract 合约项目

    合约项目中主要有两块内容：
    1、合约文件：erc20_score_contract/contracts/MyERC20.sol
      该合约可以通过配置部署到local、sepolia、base sepolia等链
      
    2、合约测试文件：erc20_score_contract/test/01_test.js
      执行该文件可以部署合约到链上，并且，该文件代码中生成了10个用户账号，通过代码中的定时任务将使用这10个账号的
      进行mint、burn、tranfer操作，模拟用户的在区块链上操作，生成区块交易数据。
      合约测试文件第七行代码为：let cycle_number = 60;
      修改cycle_number的值，可以配置交易笔数。



## 2.2 erc20_score_backend 后台项目

   
   
       
       1、.log目录用于存放日志
       2、config目录存放配置文件和数据库sql
       3、src目录存放业务代码：
          src\config 用于读取和解析配置文件
    	  src\controller 用于实现api、router、middleware
    	  src\service 用于实现业务逻辑处理
    	  src\dao  用于实现数据库连接、数据持久化和数据库模型定义
    	  src\model 用于存放数据模型，包括请求参数和返回参数
    	  src\uitil 用于存放项目工具相关方法
          main.go  项目启动文件
    	  

# 3、项目部署
## 3.1 合约部署
    `找两台电脑均部署该合约(也可分别部署到sepolia和base sepolia网络)，按照如下步骤：
    1、下载项目：git clone https://github.com/zhuruiIcarbonx/erc20_score.git
	2、进入目录：cd erc20_score/erc20_score_contract
	3、安装合约项目依赖：npm install
	4、启动合约项目hardhat节点：npx hardhat node
	5、修改模拟交易记录数量配置，默认60，可以不配置。（合约测试文件erc20_score_contract/test/01_test.js第七行代码为：let cycle_number = 60;可修改）
	6、部署合约到本地，并通过测试代码模拟交易记录：npx hardhat  --network localhost test '.\test\01_test.js'
	7、在deploy\.cache\MyERC20.json中找到合约地址："address":"0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
	
	
## 3.2 部署后台
    后台只需部署在一台电脑上即可
        1、将3.1第7步对应的两台电脑合约地址复制到erc20_score_backend\config\config.yml中，配置在blocakchain.chian1.erc20_address和blocakchain.chian2.erc20_address
        2、hardhat的chainId为31337，为了区分，分别设置blocakchain.chian1.chain_id和blocakchain.chian2.chain_id为：31337_1和31337_2
    	3、hardhat的节点默认port为8545，分别根据ip设置：blocakchain.chian1.key_url和blocakchain.chian2.key_url
		4、在config.yml中配置mysql数据库
		5、在mysql客户端执行erc20_score_backend\config\database.sql中的sql。其中最后两行insert语句中contract_address字段需要根据实际情况配置
		6、进入目录erc20_score_backend，执行命令行安装go依赖：go mod tidy
		7、启动后台项目：go run main.go
		8、此时可根据表来查询数据库数据：
			t_chain   链配置表
			t_transaction 交易表
			t_user_balance 余额表
			t_user_balance_his 余额历史表
			t_user_score 分数表
			t_user_score_his 分数历史表
    	
	
# 4、思考问题解决
考虑一个场景，如果程序错误了，或者rpc有问题，导致好几天没有计算积分。此时应该如何正确回溯？
    答：调用“历史积分计算”接口，输入参数，即可将chainId所在链从fromHour到toHour的历史积分重新计算一遍：
    接口路径：/erc20/v1/score/calculate
    请求方式： POST  
    content-type: application/json
    入参：
    {
      "chainId:"31337_1", //对应config.yml中配置的chain_id
      "fromHour:"2025-09-21 10:00:00",//需要是整点
      "toHour:"2025-09-21 15:59:59",//需要是59分59秒
    }
    出参：
    {
      "code":"0",//0：代表成功，其他：失败
      "message"："success",
      "time"："2025-09-21 18:00:01",
    }

