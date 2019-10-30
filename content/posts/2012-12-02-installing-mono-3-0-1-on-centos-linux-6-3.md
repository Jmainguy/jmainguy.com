---
title: Installing Mono 3.0.1 on Centos Linux 6.3
type: post
date: 2012-12-02T22:53:08+00:00
url: /index.php/installing-mono-3-0-1-on-centos-linux-6-3/
al2fb_facebook_link_id:
  - 55700847_10100653368910818
al2fb_facebook_link_time:
  - 2012-12-02T22:53:12+00:00
al2fb_facebook_link_picture:
  - 'avatar=http://1.gravatar.com/avatar/11181a1e72a3b6feeaec90f8a39fcc2f?s=96&amp;d=http%3A%2F%2F1.gravatar.com%2Favatar%2Fad516503a11cd5ca435acc9bb6523536%3Fs%3D96&amp;r=G'
categories:
  - General

---
So you want to run windows .exe&#8217;s on your shiny Linux server, time to catch [mono][1].
  
To install mono first install your dependencies.

<pre lang="shell" prompt="$">yum install bison gettext glib2 freetype fontconfig libpng libpng-devel libX11 libX11-devel glib2-devel libgdi* libexif glibc-devel urw-fonts java unzip gcc gcc-c++ automake autoconf libtool make bzip2 wget</pre>

<pre lang="shell" prompt="$">cd /usr/local/src 
</pre>

<pre lang="shell" prompt="$">wget http://download.mono-project.com/sources/mono/mono-3.0.1.tar.bz2
</pre>

<pre lang="shell" prompt="$">tar jxf mono-3.0.1.tar.bz2
</pre>

<pre lang="shell" prompt="$">cd mono-3.0.1
</pre>

<pre lang="shell" prompt="$">./configure --prefix=/opt/mono
</pre>

<pre lang="shell" prompt="$">make && make install
</pre>

Congrats, you now have Mono installed and can run some Windows .exe&#8217;s such as Terraria or Minecraft Classic.

 [1]: http://www.mono-project.com/Main_Page