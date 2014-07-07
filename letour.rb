require "rubygems"
require "bundler/setup"

# in stdlib
require "open-uri"
require "json"

# in the bundle
require "sinatra"

module TourHighlightsLinks
  TITLE_REGEX = /tour de france.+(stage \d+|prologue).+highlights/i

  def self.get_links
    get_links_from_json
  end

private

  def self.get_links_from_json
    feed_url = "http://www.sbs.com.au/api/video_feed/f/Bgtm9B/sbs-section-clips?form=json&byCategories=Sport%2FCycling&q=road&byCustomValue=%7BuseType%7D%7BHighlights%7D&range=1-100&defaultThumbnailAssetType=Thumbnail"

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

