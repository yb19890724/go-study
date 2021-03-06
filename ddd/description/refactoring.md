#### 重构

要避免重复建设，就要理解中台的理念和思想。前面说了“`中台是企业级能力复用平台`”，`“复用”`用白话说就是重复使用，就是要避免重复造轮子的事情。

####### 中台的设计思想与“`高内聚、低耦合`”的设计原则是高度一致的。

高内聚是把相关的业务行为聚集在一起，把不相关的行为放在其它地方，如果你要修改某个业务行为，只需要修改一处。
对了！中台就是要这样做，按照“高内聚、松耦合”的原则，实现企业级的能力复用！


###### 如何构建中台业务模型？

自顶向下和自底向上的策略。具体采用哪种策略，你需要结合公司的下面我就来介绍一下这两种策略。

1. 自顶向下的策略

这种策略是先做顶层设计，从最高领域逐级分解为中台，分别建立领域模型，根据业务属性分为通用中台或核心中台。
领域建模过程主要基于业务现状，暂时不考虑系统现状。`自顶向下的策略适用于全新的应用系统建设，或旧系统推倒重建的情况`。

由于这种策略不必受限于现有系统，你可以用 DDD 领域逐级分解的领域建模方法。
从下面这张图我们可以看出它的主要步骤：

- 第一步是将领域分解为子域，子域可以分为核心域、通用域和支撑域；
- 第二步是对子域建模，划分领域边界，建立领域模型和限界上下文；
- 第三步则是根据限界上下文进行微服务设计。

![Image text](https://static001.geekbang.org/resource/image/e6/da/e665d85381a9b2c599555cac6a06deda.jpg)

2. 自底向上的策略

这种策略是基于业务和系统现状完成领域建模。

> 自底向上策略适用于遗留系统业务模型的演进式重构。

- 首先分别完成系统所在业务域的领域建模；
- 然后对齐业务域，找出具有同类或相似业务功能的领域模型，对比分析领域模型的差异，重组领域对象，重构领域模型。这个过程会沉淀公共和复用的业务能力，会将分散的业务模型整合。


###### 第一步：锁定系统所在业务域，构建领域模型。

锁定系统所在的业务域，采用事件风暴，找出领域对象，构建聚合，划分限界上下文，建立领域模型。
看一下下面这张图，我们选取了传统核心应用的用户、客户、传统收付和承保四个业务域以及互联网电商业务域，共计五个业务域来完成领域建模。


在这些领域模型的清单里，我们可以看到二者之间有很多名称相似的领域模型。
深入分析后你会发现，这些名称相似的领域模型存在业务能力重复，或者业务职能分散（比如移动支付和传统支付）的问题。
那在构建中台业务模型时，你就需要重点关注它们，将这些不同领域模型中重复的业务能力沉淀到中台业务模型中，将分散的领域模型整合到统一的中台业务模型中，对外提供统一的共享的中台服务。

![Image text](https://static001.geekbang.org/resource/image/f5/46/f537a7a43e77212c8a85241439b2f246.jpg)

###### 第二步：对齐业务域，构建中台业务模型。

首先我们可以将传统核心的领域模型作为主领域模型，将互联网电商领域模型作为辅助模型来构建中台业务模型。
然后再将互联网电商中重复的能力沉淀到传统核心的领域模型中，只保留自己的个性能力，比如订单。
中台业务建模时，既要关注领域模型的完备性，也要关注不同渠道敏捷响应市场的要求。

![Image text](https://static001.geekbang.org/resource/image/25/1d/25cd1e7fe14bfa22a752c1b184b9c91d.jpg)


###### 中台业务模型的构建过程

> 构建多业务域的中台业务模型的过程，就是找出同一业务域内所有同类业务的领域模型，对比分析域内领域模型和聚合的差异和共同点，打破原有的模型，完成新的中台业务模型重组或归并的过程。

![Image text](https://static001.geekbang.org/resource/image/fb/70/fb11e6941fc471c734d0b85c25cc5370.jpg)

###### 第三步：中台归类，根据领域模型设计微服务。

完成中台业务建模后，我们就有了下面这张图。
从这张图中我们可以看到总共构建了多少个中台，中台下面有哪些领域模型，哪些中台是通用中台，哪些中台是核心中台，中台的基本信息等等，都一目了然。
你根据中台下的领域模型就可以设计微服务了。

![Image text](https://static001.geekbang.org/resource/image/a8/c5/a88e9695c7198a1f88f537564ada0bc5.jpg)

###### 重构过程中的领域对象

上面主要是从聚合的角度来描述中台业务模型的重组，是相对高阶的业务模块的重构。

业务模型重构和聚合重组，往往会带来领域对象和业务行为的变化。
下面我带你了解一下，在领域模型重组过程中，发生在更底层的领域对象的活动。

传统核心客户领域模型重构之前，包含个人、团体和评级三个聚合，每个聚合内部都有自己的聚合根、实体、方法和领域服务等。

![Image text](https://static001.geekbang.org/resource/image/ae/3c/ae33bc5c0cda28740363e39edbc1e53c.jpg)

互联网电商客户领域模型重构前包含个人和积分两个聚合，每个聚合包含了自己的领域对象、方法和领域服务等。

![Image text](https://static001.geekbang.org/resource/image/d0/7d/d0f8fb06797a5983c7fd00d59d8be57d.jpg)


总结：
其实呢，中台业务模型的重构过程，也是微服务架构演进的过程。业务边界即微服务边界，业务边界做好了，微服务的边界自然就会很好。