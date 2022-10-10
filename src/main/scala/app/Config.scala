package app

/**
 * @param mode 動作モード
 * @param baseDirectoryPath ベースディレクトリの位置
 * @param templatePathOpt テンプレートファイル
 * @param startsWithOpt アーカイブ対象
 */
case class Config(
    mode: String = "",
    baseDirectoryPath: String = "",
    templatePathOpt: Option[String] = None,
    startsWithOpt: Option[String] = None
)
