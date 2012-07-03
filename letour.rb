require "rubygems"
require "bundler/setup"

# in stdlib
require "open-uri"

# in the bundle
require "nokogiri"
require "sinatra"

module TourHighlightsLinks
  def self.get_links
    sbs_domain = "http://www.sbs.com.au"
    videos_url = sbs_domain + "/cyclingcentral/videos"

    videos_page_doc = Nokogiri::HTML(open(videos_url))

    highlights_title_regexp = /tour de france stage \d+ extended highlights/i

    links = videos_page_doc.css("a").select { |l| 
        !l.attributes["rel"] && # ignore anything with a rel (it's for fancy javascript)
        l.attributes["title"] && l.attributes["title"].value =~ highlights_title_regexp 
      }.collect { |l| 
        [l.attributes["title"].value, sbs_domain + l.attributes["href"].value]
      }
  end
end

get '/' do
  TourHighlightsLinks.get_links.to_s
end

