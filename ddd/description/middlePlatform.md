#### 中台到底是什么？

    中台是一个基础的理念和架构，我们要把所有的基础服务用中台的思路建设，进行联通，共同支持上端的业务。
    业务中台更多的是支持在线业务，数据中台提供了基础数据处理能力和很多的数据产品给所有业务方去用。
    业务中台、数据中台、算法中台等等一起提供对上层业务的支撑。”
    
    
> 中台的关键词：共享、联通、融合和创新。

中台首先体现的是一种企业级的能力，它提供的是一套企业级的整体解决方案，解决小到企业、集团，大到生态圈的能力共享、联通和融合问题，支持业务和商业模式创新。通过平台联通和数据融合为用户提供一致的体验，更敏捷地支撑前台一线业务。


中台来源于平台，但中台和平台相比，它更多体现的是一种理念的转变，它主要体现在这三个关键能力上：
- 对前台业务的快速响应能力；
- 企业级复用能力；
- 从前台、中台到后台的设计、研发、页面操作、流融合能力。程服务和数据的无缝联通、融合能力。


###### 如何实现前中后台的协同？

传统企业的早期系统有不少是基于业务领域或组织架构来建设的，每个系统都有自己的前端，相互独立，用户操作是竖井式，需要登录多个系统才能完成完整的业务流程。

![Image text](https://static001.geekbang.org/resource/image/76/91/76c677ccc83912dbc4d09d62c259b391.jpg)

中台后的前台建设要有一套综合考虑业务边界、流程和平台的整体解决方案，以实现各不同中台前端操作、流程和界面的联通、融合。
不管后端有多少个中台，前端用户感受到的就是只有一个前台。

在前台设计中我们可以借鉴微前端的设计思想，在企业内不仅实现前端解耦和复用，还可以根据核心链路和业务流程，通过对微前端页面的动态组合和流程编排，实现前台业务的融合。

在前台设计中我们可以借鉴微前端的设计思想，在企业内不仅实现前端解耦和复用，还可以根据核心链路和业务流程，通过对微前端页面的动态组合和流程编排，实现前台业务的融合。前端页面可以很自然地融合到不同的终端和渠道应用核心业务链路中，实现前端页面、流程和功能复用。


#### 中台

传统企业的核心业务大多是基于集中式架构开发的，而单体系统存在扩展性和弹性伸缩能力差的问题，因此无法适应忽高忽低的互联网业务场景。
而数据类应用也多数通过 ETL 工具抽取数据实现数据建模、统计和报表分析功能，但由于数据时效和融合能力不够，再加上传统数据类应用本来就不是为前端而生的，因此难以快速响应前端一线业务。


> 业务中台的建设可采用领域驱动设计方法，通过领域建模，将可复用的公共能力从各个单体剥离，沉淀并组合，采用微服务架构模式，建设成为可共享的通用能力中台。

同样的，我们可以将核心能力用微服务架构模式，建设成为可面向不同渠道和场景的可复用的核心能力中台。 业务中台向前台、第三方和其它中台提供 API 服务，实现通用能力和核心能力的复用。