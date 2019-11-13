# 写在前面

- 项目地址：[Github](https://github.com/hsyjkjkl/ServiceComputingonCloud/tree/master/web)

- 资料来源于老师的课程网页：[传送门](https://pmlpml.github.io/ServiceComputingOnCloud/ex-cloudgo-inout)

- 参考资料：

	- 老师提供的资料：[golang web 服务器 request 与 response 处理](https://blog.csdn.net/pmlpml/article/details/78539261)
  
	- 表单的处理方式：[build-web-application-with-golang](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/04.0.md)

# 实验步骤

## 1. 静态文件服务、模板输出

首先构建一个main文件，用来设置网站搭建的端口号（IP是本机地址127.0.0.1），这里会使用到pflag的知识，通过命令行参数的方式将设置端口号。

```go
package main

import (
	"os"

	"github.com/hsyjkjkl/web/service"
	flag "github.com/spf13/pflag"
)

const (
	PORT string = "8000"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	pPort := flag.StringP("port", "p", PORT, "PORT for httpd listening")
	flag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}

	server := service.NewServer()
	server.Run(":" + port)
}

```

然后调用server的函数来构建这个网站，包括html文件渲染，静态目录的支持：

```go
package service

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
		Directory:  "templates",
		Extensions: []string{".html"},
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	webRoot := os.Getenv("WEBROOT")
	if len(webRoot) == 0 {
		if root, err := os.Getwd(); err != nil {
			panic("Could not retrive working directory")
		} else {
			webRoot = root
		}
	}

	mx.PathPrefix("/").Handler(http.StripPrefix("", http.FileServer(http.Dir(webRoot+"/assets/"))))

}
```

其中formatter使用了render来指定了模板的目录，模板文件的扩展名，以便直接输出一个html文件。这个可以作为访问的主页来进行展示，这样就相当于制定了templates目录下的html文件作为渲染的目标，处理路由的时候就不必再特殊处理这个路径。
实现静态文件服务最关键的就是前缀处理，也就是如何从url映射到相应的文件目录中，就是一行代码：

```go
mx.PathPrefix("/").Handler(http.StripPrefix("", http.FileServer(http.Dir(webRoot+"/assets/"))))
```

将url中/之后的后缀，映射到WEBROOT/assets/中去，这里WEBROOT为空，没有配置。所以就是当前目录下的assets文件夹的文件都可以利用url来访问。
访问`css/main.css`:

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113163324422.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0pLSktMMQ==,size_16,color_FFFFFF,t_70)

访问`js/hello.js`:

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113163447841.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0pLSktMMQ==,size_16,color_FFFFFF,t_70)

当我们在当前根目录下没有添加html文件的时候，就会显示出这个文件目录结构：

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113163906171.png)

而如果添加了index.html，这个文件内容就会自动渲染为网页，因为我们调用的`http.FileServer`这个方法的实现里面，使用到了`serverFile`会判断是否存在index.html这个文件，如果存在，就会将其重定向到./，也就是当前的url中，所以尽管访问的url是`XXX/`，但是实际上访问的是`XXX/index.html`，所以就会看到html页面的内容而不是文件目录本身。

如果想要直接在主页显示templates目录下的文件夹，就可以选择添加一个Handler：

```go
mx.HandleFunc("/", templateHandler(formatter)).Methods("GET")
```

注意到在这里利用了，之前定义的formatter，formatter里面记录的文件模板目录的路径，以及文件的扩展名，所以我们在实现`templateHandler`这个函数的时候，只需要指定文件名字也就是index就可以了。

```go
package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func templateHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		formatter.HTML(w, http.StatusOK, "index", struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}{ID: "17343038", Content: "Hello Web!"})
	}
}

```

这里对index.html中定义的ID和Content两个属性进行赋值操作，指定为string类型，是json格式的。

最后直接访问`localhost:8000`得到的结果就是`templates/index.html`的页面内容：

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113170426503.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0pLSktMMQ==,size_16,color_FFFFFF,t_70)

红线标出的地方就是使用`templateHandler`这个方法，将html中的属性赋上了值。

## 2. 支持简单 js 访问

首先需要编写另一个url处理函数，使得js能够请求相应的数据，并且添加到html页面上作为验证。

```go
package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func jsonHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}{ID: "8675309", Content: "Hello from Go!"})
	}
}

```

上边直接从老师上课的资料中复制过来，更改了函数名字，方便辨认。可以看到这里使用的是`formatter.JSON`这个方法，返回一个JSON格式的数据，key值与之前一样，但是这次不是通过formatter来添加到html上，而是通过js来访问获得数据，再添加到页面上。由于主页已经使用了`formatter.HTML`的方式来实现，所以需要另外再assets中添加一个html文档，这里也是直接用的老师上课内容的代码：

```html
<html>
<head>
  <link rel="stylesheet" href="../css/main.css"/>
  <script src="http://code.jquery.com/jquery-latest.js"></script>
  <script src="../js/hello.js"></script>
</head>
<body>
	<div>
		<p>Sample Go Web Application!!</p>
		<p class="greeting-id">The ID is </p>
		<p class="greeting-content">The content is </p>
	</div>
</body>
</html>
```

JS的异步请求以及数据添加实现，使用到了JQuery来实现：

```javascript
$(document).ready(function() {
    $.ajax({
        url: "/api/test"
    }).then(function(data) {
       $('.greeting-id').append(data.id);
       $('.greeting-content').append(data.content);
    });
});
```

首先获取/api/test的请求数据，然后将数据中的id和content解析出来添加到html页面中相应的类的标签内容中去。

然后在server.go中添加相应的url处理函数：

```go
mx.HandleFunc("/api/test", jsonHandler(formatter)).Methods("GET")
```

此时访问`localhost:8000/api/test`会看到JSON格式的数据：

![在这里插入图片描述](https://img-blog.csdnimg.cn/2019111317154952.png)

访问`localhost:8000/html/`（存放新的html文件的目录），就会自动将里面的`index.html`文件渲染出来，并且其中引用的js文件脚本就会被执行：

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113171808837.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0pLSktMMQ==,size_16,color_FFFFFF,t_70)

结果就是JS将请求得到的数据添加到html的内容中去。

## 3. 表单数据处理

html中添加表单的标签，并且添加一个提交按钮，定向到url: `/submit`中去，方法为POST：

```html
<form action="/submit" method="post">
     <label for="username">姓名:</label>
     <input type="text" name="username" placeholder="该怎么称呼你呢~">
     <br />
     <label for="password">学校:</label>
     <input type="text" name="university" placeholder="输入你的学校">
     <br />
     <input type="submit" value="提交">
</form>
```

然后试着点击在页面中点击提交的按钮，发现会跳转到`localhost:8000/submit`：

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113172947886.png)

而我们要处理的就是当跳转到这个页面的时候需要怎么处理数据，以及如何输出一个表格。同样的，先添加一个处理函数：

```go
mx.HandleFunc("/submit", submitHandler(formatter)).Methods("POST")
```

然后实现submitHandler：

```go
package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func submitHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		formatter.HTML(w, http.StatusOK, "form", struct {
			NAME       string `json:"name"`
			UNIVERSITY string `json:"university"`
		}{NAME: req.FormValue("username"), UNIVERSITY: req.FormValue("university")})
	}
}
```

这里与之前稍微有点不同的是，由于需要对表单提交的数据进行处理，所以需要对req进行一个解析，这里直接调用`req.FormValue`会自动地解析Form的数据，并且获取表单某个名字的输入框的提交值。

注意如果有多个同名的输入框，该函数只会获取到第一个提交的值。如果需要多个一并获取，就需要用到两个函数， `req.ParseForm`解析Form 和 `req.Get("XXX")`或者`req.Form["XXX"]`取得某个name的输入框的值。这里返回的是一个slice，所以如果需要取出特定的元素，还需要加上下标。（参考：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/04.1.md）


最后实现输出表格的html模板：

```html
<!DOCTYPE html>
<html>
<head>
    <style type="text/css">
        table{
            text-align: center;
        }
    </style>
</head>

<body>
    <table border="1" cellpadding="10">
        <caption>Hello, {{.NAME}}</caption>
        <tr>
            <th>姓名</th>
            <th>学校</th>
        </tr>
        <tr>
            <td>{{.NAME}}</td>
            <td>{{.UNIVERSITY}}</td>
        </tr>
    </table>
</body>
</html>
```

以上工作就会在submit按钮提交之后，调用handler处理表单数据，然后对`templates/form.html`进行属性值填充，最后渲染出来。结果：

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113180414494.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0pLSktMMQ==,size_16,color_FFFFFF,t_70)

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113180429444.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0pLSktMMQ==,size_16,color_FFFFFF,t_70)

## 4. 对 /unknown 给出开发中的提示

跟老师上课的练习一样：

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113180837482.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0pLSktMMQ==,size_16,color_FFFFFF,t_70)

只需要稍微更改下url和文字内容就可以了。

首先添加一个url处理：

```go
mx.HandleFunc("/unknown", notImplemented).Methods("GET")
```

再实现相应的处理函数：

```go
func notImplemented(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "501 Not Implemented\nWe are still working on it. Please wait!", 501)
}
```

也就调用`http.Error`的事情，因为它的`NotFound`函数也是这样实现的。

结果：

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113181757784.png)

本次实验暂时到这里！

之后补充curl和ab测试的内容！
