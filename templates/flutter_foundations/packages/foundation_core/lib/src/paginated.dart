final class Paginated<T> {
  const Paginated({
    required this.items,
    required this.page,
    required this.perPage,
    required this.total,
  });
  final List<T> items;
  final int page;
  final int perPage;
  final int total;
  bool get hasNext => page * perPage < total;
}
