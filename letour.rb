require "rubygems"
require "bundler/setup"

# in stdlib
require "open-uri"

# in the bundle
require "nokogiri"
require "pry"

sbs_domain = "http://www.sbs.com.au"
videos_url = sbs_domain + "/cyclingcentral/videos"

videos_page_doc = Nokogiri::HTML(open(videos_url))

links = videos_page_doc.css("a").select { |l| l.attributes["title"] && !l.attributes["rel"] && l.attributes["title"].value =~ /tour de france stage \d+ extended highlights/i }.collect { |l| [l.attributes["title"].value, sbs_domain + l.attributes["href"].value]}

binding.pry
