require "rubygems"
require "bundler/setup"

# in stdlib
require "open-uri"
require "json"

# in the bundle
require "nokogiri"
require "sinatra"

module TourHighlightsLinks
  TITLE_REGEX = /tour de france.+highlights.+(stage \d+|prologue)/i

  def self.get_links
    get_links_from_json
  end

private
  def self.get_links_from_html
    sbs_domain = "http://www.sbs.com.au"
    videos_url = sbs_domain + "/cyclingcentral/videos"

    videos_page_doc = Nokogiri::HTML(open(videos_url))

    videos_page_doc.css("a").select { |l| 
      !l.attributes["rel"] && # ignore anything with a rel (it's for fancy javascript)
      l.attributes["title"] && l.attributes["title"].value =~ TITLE_REGEX 
    }.collect { |l| 
      [l.attributes["title"].value, sbs_domain + l.attributes["href"].value]
    }
  end

  def self.get_links_from_json
    feed_url = "http://www.sbs.com.au/api/video_feed/f/Bgtm9B/sbs-section-sbstv/?range=1-100&byCategories=Sport/Cycling&form=json&defaultThumbnailAssetType=Thumbnail"

    video_json = JSON.load(open(feed_url))

    video_json["entries"].select { |entry| 
      entry["title"] =~ TITLE_REGEX
    }.collect { |entry|
      HighlightsVideo.new(entry)
    }.collect { |highlight|
      [highlight.title,
       highlight.video_url]
    }
  end

  class HighlightsVideo
    def initialize(json_entry)
      @json_entry = json_entry
    end

    def title
      @json_entry["title"]
    end

    def video_url
      download_url_key = "plfile$downloadUrl"
      preferred_bit_rate = "1500K"

      @json_entry["media$content"].select { |content| 
        content[download_url_key].include?(preferred_bit_rate)
      }.first[download_url_key]
    end
  end
end

get '/' do
  erb :index, :locals => { :video_links => TourHighlightsLinks.get_links }
end

