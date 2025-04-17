<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet version="3.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
  xmlns:atom="http://www.w3.org/2005/Atom" xmlns:base="http://purl.org/atompub/base/1.0/" xmlns:str="https://tduyng.com">
  <xsl:output method="html" version="1.0" encoding="UTF-8" indent="yes"/>
  <xsl:template match="/">
    <html xmlns="http://www.w3.org/1999/xhtml" lang="en">
      <head>
        <title>
          <xsl:value-of select="/atom:feed/atom:title"/> • Feed
        </title>
        <meta charset="utf-8"/>
        <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
        <meta name="viewport" content="width=device-width, initial-scale=1"/>
        <link rel="stylesheet" href="{/atom:feed/atom:link[@rel='extra-stylesheet']/@href}" />

      </head>
      <body dir="auto">
        <div>
          <main>
            <div>
                <xsl:choose>
                    <xsl:when test="contains(/atom:feed/str:translations/str:about_feeds, 'About Feeds')">
                        <xsl:value-of select="substring-before(/atom:feed/str:translations/str:about_feeds, 'About Feeds')"/>
                        <a href="https://aboutfeeds.com/" target="_blank">About Feeds</a>
                        <xsl:value-of select="substring-after(/atom:feed/str:translations/str:about_feeds, 'About Feeds')"/>
                    </xsl:when>
                    <xsl:otherwise>
                        <xsl:value-of select="/atom:feed/str:translations/str:about_feeds"/>
                    </xsl:otherwise>
                </xsl:choose>
            </div>
            <section >
            <div >
                <xsl:value-of select="/atom:feed/atom:title"/>
            </div>
            <p>
              <xsl:value-of select="/atom:feed/atom:subtitle"/>
            </p>
            <a >
              <xsl:attribute name="href">
                <xsl:value-of select="/atom:feed/@xml:base"/>
              </xsl:attribute>
            <xsl:value-of select="/atom:feed/str:translations/str:visit_the_site" />&#160;<span class="arrow">→</span></a><p></p>
            </section>

            <div >
            <h1><xsl:value-of select="/atom:feed/str:translations/str:recent_posts" /></h1>
            </div>
            <div>
              <xsl:for-each select="/atom:feed/atom:entry">
                <section >
                  <ul >
                    <li ><xsl:value-of select="substring(atom:published, 0, 11)"/></li>
                  </ul>
                  <div >
                    <div >
                      <a>
                        <xsl:attribute name="href">
                          <xsl:value-of select="atom:link/@href"/>
                        </xsl:attribute>
                        <xsl:value-of select="atom:title"/>
                      </a>
                    </div>
                    <div>
                      <xsl:value-of select="atom:summary"/>
                    </div>
                    <a  href="">
                      <xsl:attribute name="href">
                        <xsl:value-of select="atom:link/@href"/>
                      </xsl:attribute>
                    </a>
                  </div>
                </section>
              </xsl:for-each>
            </div>
          </main>
        </div>
      </body>
    </html>
  </xsl:template>
</xsl:stylesheet>
