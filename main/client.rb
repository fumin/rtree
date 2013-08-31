require 'json'
require 'socket'

class Client
  def initialize host, port
    @sock = TCPSocket.new host, port
  end

  def rtree_insert key, member, point, lengths
    b = {method: "Store.RtreeInsert",
         params: [{"Key"    => key,
                   "Member" => member,
                   "Where"  => {"Point" => point, "Lengths" => lengths}}],
         id:     rand(1000)}
    @sock.write(JSON.dump(b))
    resp = JSON.load(@sock.readline)
    resp["result"]["Member"]
  end

  def rtree_delete key, member
    b = {method: "Store.RtreeDelete",
         params: [{"Key" => key, "Member" => member}],
         id: rand(1000)}
    @sock.write(JSON.dump(b))
    resp = JSON.load(@sock.readline)
    resp["result"]
  end

  def rtree_nearest_neighbors key, k, point
    b = {method: "Store.RtreeNearestNeighbors",
         params: [{"Key" => key, "K" => k, "Point" => point}],
         id: rand(1000)}
    @sock.write(JSON.dump(b))
    resp = JSON.load(@sock.readline)
    resp["result"]["Members"]
  end
end

if __FILE__ == $0
  c = Client.new 'localhost', 6389
  c.rtree_insert "test", "a", [0, 0], [3.5, 4.2]
  c.rtree_insert "test", "b", [2, 2], [3.1, 2.7]
  c.rtree_insert "test", "c", [100, 100], [1.1, 1.2]
  p c.rtree_nearest_neighbors "test", 5, [3.4, 4.201]

  c.rtree_delete "test", "b"
  neighbors = c.rtree_nearest_neighbors("test", 5, [3.4, 4.201])
  puts "After deleting b: #{neighbors}"

  c.rtree_insert "test", "c", [1.2, 2], [3, 3.1]
  neighbors = c.rtree_nearest_neighbors("test", 5, [3.4, 4.201])
  puts "After making c closer: #{neighbors}"
end
