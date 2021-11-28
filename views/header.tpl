<div class="layui-header layui-bg-cyan">
      <div class="layui-logo layui-hide-xs layui-bg-cyan">{{.AppName}}</div>
      <!-- 头部区域（可配合layui 已有的水平导航） -->
      <ul class="layui-nav layui-layout-left">

        <li class="layui-nav-item layui-show-xs-inline-block layui-hide-sm" lay-header-event="menuLeft">
          <i class="layui-icon layui-icon-spread-left"></i>
        </li>
        <li class="layui-nav-item">
          <a href="javascript:;">文件处理</a>
          <dl class="layui-nav-child">
            {{str2html .Nav}} 
          </dl>
        </li>
      </ul>

      <!-- 个人信息 -->
      <ul class="layui-nav layui-layout-right">
        <li class="layui-nav-item layui-hide layui-show-md-inline-block">
          <a href="javascript:;">
            <img src="//tva1.sinaimg.cn/crop.0.0.118.118.180/5db11ff4gw1e77d3nqrv8j203b03cweg.jpg"
              class="layui-nav-img">
            {{.Uid}}
          </a>
          <dl class="layui-nav-child">
            <dd><a href="">Your Profile</a></dd>
            <dd><a href="">Settings</a></dd>
            <dd><a href="">Sign out</a></dd>
          </dl>
        </li>
        <li class="layui-nav-item" lay-header-event="menuRight" lay-unselect>
          <a href="javascript:;">
            <i class="layui-icon layui-icon-more-vertical"></i>
          </a>
        </li>
      </ul>
    </div>